package pagemanager

import (
	"database/sql"
	"net/http"
	"sync"

	"github.com/bokwoon95/pagemanager/internal/tables"
	"github.com/bokwoon95/pagemanager/sq"
	"github.com/bokwoon95/pagemanager/sq/ddl"
)

type Plugin interface {
	DefaultName() string
	URL() string
	Version() string
	Capabilities() []string
	Roles() map[string][]string
	Setup(cfg *Config) error
	Middleware() func(http.Handler) http.Handler
	HandlerFuncs() map[string]http.HandlerFunc
}

var (
	pluginsMu    sync.RWMutex
	plugins      []Plugin
	pluginsIndex = make(map[[2]string]int)
)

func Register(plugin Plugin) {
	if plugin == nil {
		panic("pagemanager: Register plugin is nil")
	}
	pluginsMu.Lock()
	defer pluginsMu.Unlock()
	pluginURL := plugin.URL()
	if pluginURL == "" {
		panic("pagemanager: Register plugin URL is empty")
	}
	pluginVersion := plugin.Version()
	if _, dup := pluginsIndex[[2]string{pluginURL, pluginVersion}]; dup {
		pluginURLVersion := pluginURL + "@" + pluginVersion
		if pluginVersion == "" {
			pluginURLVersion = pluginURLVersion[:len(pluginURLVersion)-1]
		}
		panic("pagemanager: Register called twice for " + pluginURLVersion)
	}
	plugins = append(plugins, plugin)
	pluginsIndex[[2]string{pluginURL, pluginVersion}] = len(plugins) - 1
}

type pmPlugin struct{}

func (p *pmPlugin) DefaultName() string { return "pagemanager" }

func (p *pmPlugin) URL() string { return "github.com/pagemanager/pagemanager" }

func (p *pmPlugin) Version() string { return "" }

var capabilities = []string{
	// site-level capabilities
	"administrate-sites",
	"administrate-roles",
	"administrate-tags",
	"assign-roles",
	// url-level capabilities
	"url.administrate",
	"url.edit-all-entries",
	"url.edit-all-handlers",
	"url.edit-all-handler-configs",
	"url.edit-handler",
	"url.edit-handler-config",
}

func (p *pmPlugin) Capabilities() []string {
	return []string{
		// site-level capabilities
		"administrate-sites",
		"administrate-roles",
		"administrate-tags",
		"assign-roles",
		// url-level capabilities
		"url.administrate",
		"url.edit-all-entries",
		"url.edit-all-handlers",
		"url.edit-all-handler-configs",
		"url.edit-handler",
		"url.edit-handler-config",
	}
}

var roles = map[string][]string{
	"superadmin": {"administrate-sites", "administrate-roles", "administrate-tags"},
	"admin":      {"assign-roles"},
}

func (p *pmPlugin) Roles() map[string][]string {
	return map[string][]string{
		"superadmin": {"administrate-sites", "administrate-roles", "administrate-tags"},
		"admin":      {"assign-roles"},
	}
}

func (p *pmPlugin) Setup(cfg *Config) error {
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
			return err
		}
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
	db, err := sql.Open(driverName, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	gotMetadata, err := ddl.NewDatabaseMetadata(cfg.DatabaseDialect, ddl.WithDB(db, nil))
	if err != nil {
		return err
	}
	_ = gotMetadata
	wantMetadata, err := ddl.NewDatabaseMetadata(cfg.DatabaseDialect, ddl.WithTables(tbls...))
	if err != nil {
		return err
	}
	_ = wantMetadata
	// m, err := ddl.Migrate(ddl.DropExtraneous, gotMetadata, ddl.DatabaseMetadata{})
	m, err := ddl.Migrate(ddl.CreateMissing, gotMetadata, wantMetadata)
	if err != nil {
		return err
	}
	_ = m
	// err = m.WriteSQL(os.Stdout)
	// if err != nil {
	// 	return err
	// }
	err = m.Exec(db)
	if err != nil {
		return err
	}
	return nil
}
