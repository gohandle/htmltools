package rblog

import (
	"log"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

// Conf configures logs
type Conf struct {
	DevMode bool `env:"RB_LOG_DEV_MODE"`
}

// ParseConf parses the configuration from the env
func ParseConf() (cfg Conf, err error) {
	return cfg, env.Parse(&cfg)
}

// New creates our logger
func New(cfg Conf) (*zap.Logger, error) {
	if cfg.DevMode {
		return zap.NewDevelopment()
	}

	return zap.NewProduction()
}

// NewStd outputs a standard lib logger. Unfortunately we cannot use regular di here since
// it isn't bootstrapped yet.
func NewStd() *log.Logger {
	var cfg Conf
	l, _ := New(cfg)
	return zap.NewStdLog(l)
}
