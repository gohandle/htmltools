package i18n

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

// Conf configures the i18n setup
type Conf struct {
	DefaultLanguage string   `env:"RB_I18N_DEFAULT_LANGUAGE" envDefault:"en-US"`
	Dir             string   `env:"RB_I18N_DIR"`
	FileNames       []string `env:"RB_I18N_FILE_NAMES" envSeparator:":"`
}

// ParseConf parses the env
func ParseConf() (cfg Conf, err error) {
	return cfg, env.Parse(&cfg)
}

// MessageFiles are translations files on a filesystem
type MessageFiles fs.FS

// FromEmbed creates a translation message files fs from an embedded filesystem
func FromEmbed(efs embed.FS) MessageFiles { return efs }

// FromDir provides translation files from an actual directory
func FromDir(logs *zap.Logger, cfg Conf) MessageFiles {
	logs.Info("configure message files from fs dir", zap.String("dir", cfg.Dir))
	return os.DirFS(cfg.Dir)
}

// Params are parameters for the bundle constructor
type Params struct {
	fx.In
	Files MessageFiles
}

// New inits the i18n bundle
func New(logs *zap.Logger, cfg Conf, p Params) (b *i18n.Bundle, err error) {
	dt, err := language.Parse(cfg.DefaultLanguage)
	if err != nil {
		return nil, fmt.Errorf("invalid default language: %w", err)
	}

	b = i18n.NewBundle(dt)
	for _, fn := range cfg.FileNames {
		if err = loadMessageFile(logs, p.Files, fn, b); err != nil {
			return nil, fmt.Errorf("failed to load message file from fs: %w", err)
		}
	}

	logs.Info("finished setting up i18n bundle",
		zap.String("default_language", cfg.DefaultLanguage), zap.Strings("message_files", cfg.FileNames))
	return
}

func loadMessageFile(logs *zap.Logger, fsys fs.FS, name string, b *i18n.Bundle) (err error) {
	f, err := fsys.Open(name)
	if err != nil {
		return fmt.Errorf("failed to open message file: %w", err)
	}

	defer f.Close()
	buf, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read all files: %w", err)
	}

	var mf *i18n.MessageFile
	if mf, err = b.ParseMessageFileBytes(buf, name); err != nil {
		return fmt.Errorf("failed to parse message files: %w", err)
	}

	logs.Info("parsed message file",
		zap.String("path", mf.Path), zap.String("format", mf.Format), zap.Stringer("tag", mf.Tag))
	return
}
