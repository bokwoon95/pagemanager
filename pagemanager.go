package pagemanager

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Mode int

const (
	ModeOffline Mode = iota
	ModeSinglesite
	ModeMultisite
)

type Config struct {
	Sitemode    Mode
	DatabaseURL string
	RootFS      fs.FS

	// if empty, derive from DSN
	DatabaseURL2 string
	DatabaseURL3 string
	// if nil, derive from RootFS
	TemplatesFS fs.FS
	UploadsFS   fs.FS
}

var (
	flagSecretsFile = flag.String("pm-secrets-file", "", "")
	flagSecretsEnv  = flag.Bool("pm-secrets-env", false, "")
)

// DefaultConfig looks at the environment and the flags passed to it and
// deduces the SiteMode, DSN and RootFS from it.
func DefaultConfig() (*Config, error) {
	// -pm-secrets-file -pm-secrets-env
	// contains: PM_DSN, PM_DSN2, PM_DSN3
	cfg := &Config{}
	if *flagSecretsFile != "" {
		f, err := os.Open(*flagSecretsFile)
		if err != nil {
			return nil, fmt.Errorf("opening %s: %w", *flagSecretsFile, err)
		}
		envMap, err := godotenv.Parse(f)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", *flagSecretsFile, err)
		}
		cfg.DatabaseURL = envMap["PM_DATABASE_URL"]
		cfg.DatabaseURL2 = envMap["PM_DATABASE_URL_2"]
		cfg.DatabaseURL3 = envMap["PM_DATABASE_URL_3"]
	} else if *flagSecretsEnv {
		cfg.DatabaseURL = os.Getenv("PM_DATABASE_URL")
		cfg.DatabaseURL2 = os.Getenv("PM_DATABASE_URL_2")
		cfg.DatabaseURL3 = os.Getenv("PM_DATABASE_URL_3")
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
