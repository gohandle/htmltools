package rbasset

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

// Conf configurest static file handler
type Conf struct {
	Dir              string `env:"RB_STATIC_DIR" envDefault:"."`
	StaticPathPrefix string `env:"RB_ASSET_STATIC_PATH_PREFIX" envDefault:"/static/"`
}

// ParseConf parses the env
func ParseConf() (cfg Conf, err error) {
	return cfg, env.Parse(&cfg)
}

// StaticWare is middleware that shows static files for a certain prefix
type StaticWare func(http.Handler) http.Handler

// StaticFiles provides files for the static ware
type StaticFiles fs.FS

// FromEmbed creates a template files fs from an embedded filesystem
func FromEmbed(efs embed.FS) StaticFiles { return efs }

// FromDir provides template files from an actual directory
func FromDir(logs *zap.Logger, cfg Conf) StaticFiles {
	logs.Info("configure static dir fs", zap.String("dir", cfg.Dir))
	return os.DirFS(cfg.Dir)
}

// NewStatic returns a static file handler configured on the prefix
func NewStatic(logs *zap.Logger, cfg Conf, sfs StaticFiles) StaticWare {
	sfh := http.StripPrefix(cfg.StaticPathPrefix, http.FileServer(http.FS(sfs)))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := strings.TrimPrefix(r.URL.Path, cfg.StaticPathPrefix)
			rp := strings.TrimPrefix(r.URL.RawPath, cfg.StaticPathPrefix)
			if len(p) < len(r.URL.Path) && (r.URL.RawPath == "" || len(rp) < len(r.URL.RawPath)) {
				sfh.ServeHTTP(w, r)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
