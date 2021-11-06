package pagemanager

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bokwoon95/pagemanager/internal/tables"
	"github.com/bokwoon95/pagemanager/sq"
	"github.com/bokwoon95/pagemanager/sq/ddl"
	"github.com/joho/godotenv"
)

type Mode int

const (
	ModeOffline Mode = iota
	ModeSinglesite
	ModeMultisite
)

type Config struct {
	Mode            Mode
	DatabaseDialect string
	DatabaseURL     string
	RootFS          fs.FS
	TemplatesFS     fs.FS
	UploadsFS       fs.FS
}

var (
	flagSecretsFile = flag.String("pm-secrets-file", "", "")
	flagSecretsEnv  = flag.Bool("pm-secrets-env", false, "")
	flagRootDir     = flag.String("pm-root-dir", "", "")
	flagMode        = flag.String("pm-mode", "", "")
)

// DefaultConfig returns a default config for pagemanager based on the flags
// passed to the application.
func DefaultConfig() (*Config, error) {
	cfg := &Config{}

	switch *flagMode {
	case "0", "":
		cfg.Mode = ModeOffline
	case "1", "singlesite":
		cfg.Mode = ModeSinglesite
	case "2", "multisite":
		cfg.Mode = ModeMultisite
	default:
		return nil, fmt.Errorf("unrecognized mode: %s", *flagMode)
	}

	rootDir := *flagRootDir
	if rootDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("os.UserHomeDir: %w", err)
		}
		rootDir = filepath.Join(homeDir, "pagemanager-data")
	}
	err := os.MkdirAll(rootDir, 0775)
	if err != nil {
		return nil, fmt.Errorf("os.MkDirAll %s: %w", rootDir, err)
	}
	cfg.RootFS = os.DirFS(rootDir)

	if cfg.TemplatesFS == nil {
		templatesDir := filepath.Join(rootDir, "pm-templates")
		err = os.MkdirAll(templatesDir, 0775)
		if err != nil {
			return nil, fmt.Errorf("os.MkDirAll %s: %w", templatesDir, err)
		}
		cfg.TemplatesFS = os.DirFS(templatesDir)
	}

	if cfg.UploadsFS == nil {
		uploadsDir := filepath.Join(rootDir, "pm-uploads")
		err = os.MkdirAll(uploadsDir, 0775)
		if err != nil {
			return nil, fmt.Errorf("os.MkDirAll %s: %w", uploadsDir, err)
		}
		cfg.UploadsFS = os.DirFS(uploadsDir)
	}

	if *flagSecretsFile != "" {
		envMap := make(map[string]string)
		err := func() error {
			f, err := os.Open(*flagSecretsFile)
			if err != nil {
				return fmt.Errorf("opening %s: %w", *flagSecretsFile, err)
			}
			defer f.Close()
			envMap, err = godotenv.Parse(f)
			if err != nil {
				return fmt.Errorf("parsing %s: %w", *flagSecretsFile, err)
			}
			return nil
		}()
		if err != nil {
			return nil, err
		}
		cfg.DatabaseURL = envMap["PM_DATABASE_URL"]
	} else if *flagSecretsEnv {
		cfg.DatabaseURL = os.Getenv("PM_DATABASE_URL")
	}

	if cfg.DatabaseURL == "" {
		sqliteFile := filepath.Join(rootDir, "pm-database.sqlite")
		f, err := os.OpenFile(sqliteFile, os.O_CREATE, 0644)
		if err != nil {
			return nil, fmt.Errorf("os.OpenFile %s: %w", sqliteFile, err)
		}
		f.Close()
		cfg.DatabaseDialect = "sqlite"
		cfg.DatabaseURL = sqliteFile
	}

	if strings.HasPrefix(strings.TrimSpace(cfg.DatabaseURL), "postgres") {
		cfg.DatabaseDialect = "postgres"
	}
	return cfg, nil
}

type contextKey struct{ name string }

var (
	PagemanagerContextKey = &contextKey{name: "pagemanager"}
	ctxpool               = sync.Pool{New: func() interface{} { return &Context{} }}
)

type Context struct {
	SiteMode    Mode
	SiteID      [16]byte
	UserID      [16]byte
	Username    string
	Name        string
	Langcode    string
	TildePrefix string
	Domain      string
	Subdomain   string
	URLPath     string
	RawURL      string
}

type Pagemanager struct {
	mode         Mode
	siteID       [16]byte
	dialect      string
	db           *sql.DB
	rootFS       fs.FS
	templatesFS  fs.FS
	uploadsFS    fs.FS
	locales      map[string]string // read from locales.txt or a default
	middlewares  []func(http.Handler) http.Handler
	handlerfuncs map[[2]string]http.HandlerFunc
}

func New(cfg *Config) (*Pagemanager, error) {
	pm := &Pagemanager{
		mode:         cfg.Mode,
		dialect:      cfg.DatabaseDialect,
		rootFS:       cfg.RootFS,
		templatesFS:  cfg.TemplatesFS,
		uploadsFS:    cfg.UploadsFS,
		middlewares:  make([]func(http.Handler) http.Handler, 0, len(plugins)+1),
		handlerfuncs: make(map[[2]string]http.HandlerFunc),
	}
	var driverName string
	switch cfg.DatabaseDialect {
	case "sqlite":
		driverName = "sqlite3"
	case "postgres":
		driverName = "postgres"
	default:
		driverName = "mysql"
	}
	var err error
	pm.db, err = sql.Open(driverName, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("sql.Open %s %s: %w", driverName, cfg.DatabaseURL, err)
	}
	var tbls = []sq.SchemaTable{
		&tables.PM_SITE{},
		&tables.PM_PLUGIN{},
		&tables.PM_HANDLER{},
		&tables.PM_CAPABILITY{},
		&tables.PM_ALLOWED_PLUGIN{},
		&tables.PM_DENIED_PLUGIN{},
		&tables.PM_ALLOWED_HANDLER{},
		&tables.PM_DENIED_HANDLER{},
		&tables.PM_ROLE{},
		&tables.PM_TAG{},
		&tables.PM_ROLE_CAPABILITY{},
		&tables.PM_TAG_CAPABILITY{},
		&tables.PM_TAG_OWNER{},
		&tables.PM_USER{},
		&tables.PM_USER_ROLE{},
		&tables.PM_SESSION{},
		&tables.PM_URL{},
		&tables.PM_URL_ROLE_CAPABILITY{},
		&tables.PM_URL_TAG{},
		&tables.PM_TEMPLATE_DATA{},
	}
	for _, tbl := range tbls {
		err := sq.ReflectTable(tbl, "")
		if err != nil {
			return nil, err
		}
	}
	err = ddl.AutoMigrate(pm.dialect, pm.db, ddl.CreateMissing, ddl.WithTables(tbls...))
	if err != nil {
		return nil, fmt.Errorf("ddl.AutoMigrate: %w", err)
	}
	pm.middlewares = append(pm.middlewares, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/pm-templates") {
			}
			next.ServeHTTP(w, r)
		})
	})
	pm.handlerfuncs[[2]string{"pagemanager", "url-dashboard"}] = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world!"))
	}
	for _, plugin := range plugins {
		_ = plugin
	}
	return pm, nil
}

func (pm *Pagemanager) Pagemanager(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pmctx := ctxpool.Get().(*Context)
		defer func() {
			ctxpool.Put(pmctx)
		}()
		r = r.WithContext(context.WithValue(r.Context(), PagemanagerContextKey, pmctx))
		// if we get a pm_url hit, we build it up here instead
		if next == nil {
			next = http.DefaultServeMux
		}
		for i := range pm.middlewares {
			next = pm.middlewares[len(pm.middlewares)-1-i](next)
		}
		next.ServeHTTP(w, r)
	})
}
