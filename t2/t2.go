package t2

import (
	"html/template"

	"github.com/google/uuid"
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

// pm-templates
// pm-media

// {{ themedir . "banner.jpg" }}
// /pm-templates/{{ .ThemeDir }}/banner.jpg

// /pm-templates/plainsimple/banner.jpg
// /pm-templates/github.com/bokwoon95/plainsimple/banner.jpg

// url=/about-me handler=page params={template=plainsimple/about-me.html}
// /pm-templates/0000-0000-0000-0000/plainsimple/about-me.html
// /pm-templates/plainsimple/about-me.html

// ooh, we can have a template cache keyed by the filepath. Then instead of parsing the same template multiple times we can just stitch the already-parsed templates together.

type FS struct {
	siteMode      int // 0 - offline, 1 - singlesite, 2 - multisite
	defaultSiteID uuid.UUID
}
