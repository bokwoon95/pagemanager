package t2

import (
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/tailscale/hujson"
)

type any = interface{}

type Config struct {
	ThemeDir              string         `toml:"-"`
	TemplateFiles         []string       `toml:"template_files"`
	DataFiles             []string       `toml:"data_files"`
	DataQueries           []string       `toml:"data_queries"`
	DataFunctions         []string       `toml:"data_functions"`
	ContentSecurityPolicy map[string]any `toml:"content_security_policy"`
	// TODO: need something for templates to indicate which template functions they would like to import
	// TODO: maybe each *.html file is allowed to have its own *.html.json config file. That way parent templates don't have to declare the depedencies of the child templates. But this template depedency tree, it may spiral out to become too complex to manage. Maybe intentionally keep it simple and primitive for now.
}

type Bundle struct {
	Config
	Template *template.Template
	Data     map[string]any
}

// NOTE: if serverMode=offline and the user hits the URL dashboard, the server sets a blank 'pm-session' with no content.
// when the page handler handles a request with a 'pm-session' cookie and the server is either in offline mode or the user is a site admin, it injects javascript /pm-templates/pm-edit-button.js into the page which renders an edit button for the user to click. When the server is just crawling itself with no pm-session cookie, it does not render the edit button.

// pm-templates
// pm-media

// {{ themedir . "banner.jpg" }}
// /pm-templates/{{ .ThemeDir }}/banner.jpg

// /pm-templates/edit?data=
// No URL can be called

// site-xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// /pm-templates/plainsimple/banner.jpg
// /pm-uploads/site-xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/pm-templates/plainsimple/banner.jpg
// /pm-uploads/site-xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.jpg
// /pm-uploads/xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.jpg

// /pm-templates/plainsimple/banner.jpg
// /pm-templates/github.com/bokwoon95/plainsimple/banner.jpg

// url=/about-me handler=page params={template=plainsimple/about-me.html}
// /pm-templates/0000-0000-0000-0000/plainsimple/about-me.html
// /pm-templates/plainsimple/about-me.html

// ooh, we can have a template cache keyed by the filepath. Then instead of parsing the same template multiple times we can just stitch the already-parsed templates together.

type Handler interface {
	Name() string
	ConfigSchema() *hujson.Object
	Handler(config *hujson.Object) (http.Handler, error)
}

type Middleware interface {
	Name() string
	Middleware() func(http.Handler) http.Handler // does this need config *hujson.Object?
}

type Plugin interface {
	Name() string
	Handlers() []Handler
	Middlewares() []Middleware
}

type FS struct {
	siteMode      int // 0 - offline, 1 - singlesite, 2 - multisite
	defaultSiteID uuid.UUID
	templates     fs.FS
	uploads       fs.FS
}

func (tmplsFS *FS) ServeAssets(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/pm-templates") {
			next.ServeHTTP(w, r)
			return
		}
	})
}

type PageHandler struct {
	*FS
}

func (h *PageHandler) Name() string { return "page" }

func (h *PageHandler) ConfigSchema() *hujson.Object {
	return &hujson.Object{}
}

func (h *PageHandler) Handler(config *hujson.Object) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}), nil
}

// pm.RegisterHandlers(handlers ...Handler)
// pm.RegisterPlugins()

// /about-me?edit
// /about-me?data=data/site.json

// func ServeTemplate
