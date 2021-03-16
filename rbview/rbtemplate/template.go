package rbtemplate

import (
	"embed"
	"html/template"
	"io/fs"
	"os"

	"github.com/caarlos0/env/v6"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Conf configures the view
type Conf struct {
	Dir      string   `env:"RB_VIEW_DIR" envDefault:"."`
	Patterns []string `env:"RB_VIEW_PATTERNS" envSeparator:":" envDefault:"*.html"`
	Name     string   `env:"RB_VIEW_NAME" envDefault:"root"`
}

// ParseConf parses the env
func ParseConf() (cfg Conf, err error) {
	return cfg, env.Parse(&cfg)
}

// Files are on a filesystem
type Files fs.FS

// FromEmbed creates a template files fs from an embedded filesystem
func FromEmbed(efs embed.FS) Files { return efs }

// FromDir provides template files from an actual directory
func FromDir(logs *zap.Logger, cfg Conf) Files {
	logs.Info("configure dir fs", zap.String("dir", cfg.Dir))
	return os.DirFS(cfg.Dir)
}

// Params are parameters for view construction
type Params struct {
	fx.In
	Files Files
	Funcs []template.FuncMap `group:"rb.helper"`
}

// New creates the view templates
func New(logs *zap.Logger, cfg Conf, p Params) (*template.Template, error) {
	tmpl, hnames := template.New(cfg.Name), []string{}
	for _, fm := range p.Funcs {
		tmpl = tmpl.Funcs(fm)
		for fname := range fm {
			hnames = append(hnames, fname)
		}
	}

	logs.Info("parse templates from fs",
		zap.String("name", tmpl.Name()),
		zap.Strings("patterns", cfg.Patterns),
		zap.Strings("helpers", hnames))

	return tmpl.ParseFS(p.Files, cfg.Patterns...)
}
