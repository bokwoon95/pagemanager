package pagemanager

import (
	"database/sql"
	"io/fs"
	"net/http"
)

type Config struct {
	DSN    string
	RootFS fs.FS

	// if empty, derive from DSN
	DSNv2 string
	DSNv3 string
	// if nil, derive from RootFS
	TemplatesFS fs.FS
	UploadsFS   fs.FS
}

func DefaultConfig() {
	// -pm-secrets-file -pm-secrets-env
	// contains: DSN, DSN2, DSN3
}

type Pagemanager struct {
	db          *sql.DB
	db1         *sql.DB
	db2         *sql.DB
	rootFS      fs.FS
	templatesFS fs.FS
	uploadsFS   fs.FS
}

func New(cfg *Config) (*Pagemanager, error) {
	pm := &Pagemanager{}
	return pm, nil
}

func (pm *Pagemanager) Pagemanager(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
