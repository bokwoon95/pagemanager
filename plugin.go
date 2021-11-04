package pagemanager

import (
	"net/http"
	"strings"
	"sync"
)

type Plugin interface {
	Setup(cfg *Config) error
	DefaultName() string
	URL() string
	Version() string
	Handlers() map[string]http.Handler
	Capabilities() []string
	Roles() map[string][]string
	Middleware() func(http.Handler) http.Handler
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

var _ Plugin = (*pmPlugin)(nil)

func (p *pmPlugin) Setup(cfg *Config) error {
	return nil
}

func (p *pmPlugin) DefaultName() string { return "pagemanager" }

func (p *pmPlugin) URL() string { return "github.com/pagemanager/pagemanager" }

func (p *pmPlugin) Version() string { return "" }

func (p *pmPlugin) Handlers() map[string]http.Handler {
	return map[string]http.Handler{
		"url-dashboard": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		}),
	}
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

func (p *pmPlugin) Roles() map[string][]string {
	return map[string][]string{
		"superadmin": {"administrate-sites", "administrate-roles", "administrate-tags"},
		"admin":      {"assign-roles"},
	}
}

func (p *pmPlugin) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/pm-templates") {
			}
			next.ServeHTTP(w, r)
		})
	}
}
