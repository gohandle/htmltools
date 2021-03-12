package rbview

import (
	"io/fs"
	"os"
	"text/template"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

// Conf configures the view
type Conf struct {
	Dir      string   `env:"VIEW_DIR"`
	Patterns []string `env:"VIEW_PATTERNS" envSeparator:":" envDefault:"*.html"`
}

// ParseConf parses the env
func ParseConf() (cfg Conf, err error) {
	return cfg, env.Parse(&cfg)
}

// TemplateFiles are on a filesystem
type TemplateFiles fs.FS

// FromDir provides template files from an actual directory
func FromDir(logs *zap.Logger, cfg Conf) TemplateFiles {
	logs.Info("configure template dir fs", zap.String("dir", cfg.Dir))
	return os.DirFS(cfg.Dir)
}

// New creates the view templates
func New(logs *zap.Logger, cfg Conf, tfs TemplateFiles) (*template.Template, error) {
	logs.Info("parse templates from fs", zap.Strings("patterns", cfg.Patterns))
	return template.ParseFS(tfs, cfg.Patterns...)
}
