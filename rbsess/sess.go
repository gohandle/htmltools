package rbsess

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

// Conf configures the session store
type Conf struct {
	KeyPairs            []string `env:"RB_SESS_KEY_PAIRS" envSeparator:","`
	CookiePath          string   `env:"RB_SESS_COOKIE_PATH" envDefault:"/"`
	CookieMaxAgeSeconds int      `env:"RB_SESS_COOKIE_MAX_AGE_SECONDS" envDefault:"2592000"`
	CookieDomain        string   `env:"RB_SESS_COOKIE_DOMAIN"`
	CookieHTTPOnly      bool     `env:"RB_SESS_COOKIE_HTTP_ONLY"`
	Secure              bool     `env:"RB_SESS_SECURE"`
	CookieSameSite      int      `env:"RB_SESS_COOKIE_SAME_SITE"`
}

// ParseConf parses the environment into the session configuration
func ParseConf() (cfg Conf, err error) {
	return cfg, env.Parse(&cfg)
}

// NewCookieStore inits a cookie store
func NewCookieStore(logs *zap.Logger, cfg Conf) (sessions.Store, error) {
	var pairs [][]byte
	for _, el := range cfg.KeyPairs {
		kb, err := base64.StdEncoding.DecodeString(el)
		if err != nil {
			return nil, fmt.Errorf("failed to base64 decode session key: %v", err)
		}
		pairs = append(pairs, kb)
	}

	s := sessions.NewCookieStore(pairs...)
	configureOptions(s.Options, cfg)

	logs.Info("setup session store",
		zap.Int("num_keys", len(pairs)), zap.Any("options", s.Options))
	return s, nil
}

func configureOptions(opts *sessions.Options, cfg Conf) {
	opts.Path = cfg.CookiePath
	opts.Domain = cfg.CookieDomain
	opts.MaxAge = cfg.CookieMaxAgeSeconds
	opts.HttpOnly = cfg.CookieHTTPOnly
	opts.Secure = cfg.Secure
	opts.SameSite = http.SameSite(cfg.CookieSameSite)
}
