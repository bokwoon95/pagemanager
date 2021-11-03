package pagemanager

import (
	"context"
	"database/sql"
	"io/fs"
	"net/http"
)

type Sitemode int8

const (
	SitemodeOffline Sitemode = iota
	SitemodeSinglesite
	SitemodeMultisite
)

type Config struct {
	Sitemode Sitemode
	DSN      string
	RootFS   fs.FS

	// if empty, derive from DSN
	DSNv2 string
	DSNv3 string
	// if nil, derive from RootFS
	TemplatesFS fs.FS
	UploadsFS   fs.FS
}

// DefaultConfig looks at the environment and the flags passed to it and
// deduces the SiteMode, DSN and RootFS from it.
func DefaultConfig() *Config {
	// -pm-secrets-file -pm-secrets-env
	// contains: DSN, DSN2, DSN3
	cfg := &Config{}
	return cfg
}

type Pagemanager struct {
	sitemode    Sitemode
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
