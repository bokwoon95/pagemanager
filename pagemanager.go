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
	DatabaseURL2    string
	DatabaseURL3    string
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

	templatesDir := filepath.Join(rootDir, "pm-templates")
	err = os.MkdirAll(templatesDir, 0775)
	if err != nil {
		return nil, fmt.Errorf("os.MkDirAll %s: %w", templatesDir, err)
	}
	cfg.TemplatesFS = os.DirFS(templatesDir)

	uploadsDir := filepath.Join(rootDir, "pm-uploads")
	err = os.MkdirAll(uploadsDir, 0775)
	if err != nil {
		return nil, fmt.Errorf("os.MkDirAll %s: %w", uploadsDir, err)
	}
	cfg.UploadsFS = os.DirFS(uploadsDir)

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
		cfg.DatabaseURL2 = envMap["PM_DATABASE_URL_2"]
		cfg.DatabaseURL3 = envMap["PM_DATABASE_URL_3"]
	} else if *flagSecretsEnv {
		cfg.DatabaseURL = os.Getenv("PM_DATABASE_URL")
		cfg.DatabaseURL2 = os.Getenv("PM_DATABASE_URL_2")
		cfg.DatabaseURL3 = os.Getenv("PM_DATABASE_URL_3")
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

	if cfg.DatabaseURL2 == "" {
		cfg.DatabaseURL2 = cfg.DatabaseURL
	}

	if cfg.DatabaseURL3 == "" {
		cfg.DatabaseURL3 = cfg.DatabaseURL
	}
	return cfg, nil
}

type Pagemanager struct {
	mode        Mode
	siteID      [16]byte
	db          *sql.DB
	db1         *sql.DB
	db2         *sql.DB
	rootFS      fs.FS
	templatesFS fs.FS
	uploadsFS   fs.FS
	locales     map[string]string // read from locales.txt or a default
	// Plugins
	handlers     map[[2]string]http.Handler     // plugin.handler -> http.Handler
	capabilities map[string]map[string]struct{} // plugin -> capability -> struct{}
}

func New(cfg *Config) (*Pagemanager, error) {
	pm := &Pagemanager{}
	return pm, nil
}

type URLInfo struct {
	RawURL      string
	Domain      string
	Subdomain   string
	TildePrefix string
	Langcode    string
	URLPath     string
}

type CtxKey int

const (
	CtxKeyURLInfo CtxKey = iota
	CtxKeySiteID
	CtxKeyUserID
)

func (pm *Pagemanager) Pagemanager(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// maybe stuff everything into one key for performance?
		ctx := r.Context()
		ctx = context.WithValue(ctx, CtxKeyURLInfo, URLInfo{})
		ctx = context.WithValue(ctx, CtxKeySiteID, [16]byte{})
		ctx = context.WithValue(ctx, CtxKeyUserID, [16]byte{})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
