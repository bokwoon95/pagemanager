package templates

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"path"
	"strings"

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
}

type Bundle struct {
	Config
	Template *template.Template
	Data     map[string]any
}

type FS struct {
	templates     fs.FS
	media         fs.FS
	siteTemplates fs.FS
	siteMedia     fs.FS
}

func NewFS(one, two, three, four fs.FS) *FS {
	return &FS{
		templates:     one,
		media:         two,
		siteTemplates: three,
		siteMedia:     four,
	}
}

// TODO: oh god this is more complicated than I thought. I may be checking the site-specific media and media correctly, but what about the site-specific templates?
func (fsys *FS) resolveMedia(data map[string]any, filename string) string {
	// e.g. {{ media . "banner.jpg" }}
	themeDir, _ := data["ThemeDir"].(string)
	if themeDir == "" {
		return filename
	}
	siteID, _ := data["SiteID"].(uuid.UUID)
	if siteID != uuid.Nil {
		siteIDText := siteID.String()
		if siteIDText == "" {
			return filename
		}
		path1 := path.Join(siteIDText, "pm-templates", themeDir, filename)
		_, err := fs.Stat(fsys.siteMedia, strings.TrimLeft(path1, "/"))
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return filename
		}
		if err == nil {
			// /pm-site-media/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/pm-templates/plainsimple/banner.jpg
			return path.Join("pm-site-media", path1)
		}
	}
	path2 := path.Join("pm-templates", themeDir, filename)
	_, err := fs.Stat(fsys.siteMedia, strings.TrimLeft(path2, "/"))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return filename
	}
	if err == nil {
		// /pm-media/pm-templates/plainsimple/banner.jpg
		return path.Join("pm-media", path2)
	}
	// /pm-templates/plainsimple/banner.jpg
	return path2
}

func (fsys *FS) readFile(siteID uuid.UUID, filename string) ([]byte, error) {
	if siteID != uuid.Nil {
		siteIDText := siteID.String()
		if siteIDText == "" {
			return nil, fmt.Errorf("%v is not a valid UUID", siteID)
		}
		b, err := fs.ReadFile(fsys.siteTemplates, path.Join(siteIDText, filename))
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("reading %s: %w", path.Join(siteIDText, filename), err)
		}
		if err == nil {
			return b, nil
		}
	}
	b, err := fs.ReadFile(fsys.templates, filename)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", filename, err)
	}
	return b, nil
}

func (fsys *FS) parseTemplates(siteID uuid.UUID, themeDir string, filename string, filenames ...string) (*template.Template, error) {
	b, err := fsys.readFile(siteID, filename)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New(filename).Funcs(template.FuncMap{}).Parse(string(b))
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", filename, err)
	}
	for _, name := range filenames {
		b, err := fsys.readFile(siteID, filename)
		if err != nil {
			return nil, err
		}
		_, err = tmpl.New(name).Parse(string(b))
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", name, err)
		}
	}
	return tmpl, nil
}
