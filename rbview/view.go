package rbview

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/caarlos0/env/v6"
	"github.com/gohandle/htmltools/rbview/internal/accept"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Conf configures the binder
type Conf struct {
	EncoderOrder []string `env:"RB_VIEW_ENCODER_ORDER" envSeparator:","`
}

// ParseConf parses the env
func ParseConf() (cfg Conf, err error) {
	return cfg, env.Parse(&cfg)
}

type Encoder interface {
	MIME() string
	Encode(w http.ResponseWriter, r *http.Request, v interface{}, o Options) error
}

// Params configures dependencies binding dependencies
type Params struct {
	fx.In
	Encoders []Encoder `group:"rb.encoder"`
}

type View struct {
	encs  map[string]Encoder
	order []string
}

func New(logs *zap.Logger, cfg Conf, p Params) (v *View, err error) {
	v = &View{
		encs:  make(map[string]Encoder, len(p.Encoders)),
		order: make([]string, 0, len(p.Encoders)),
	}

	for _, enc := range p.Encoders {
		if _, ok := v.encs[enc.MIME()]; ok {
			return nil, fmt.Errorf("encoder for MIME '%s' already exists", enc.MIME())
		}

		v.encs[enc.MIME()] = enc
		v.order = append(v.order, enc.MIME())
	}

	sort.Strings(v.order)
	if len(cfg.EncoderOrder) > 0 {
		v.order = cfg.EncoderOrder
	}

	if len(v.order) != len(v.encs) {
		return nil, fmt.Errorf("encoder order specified %d encoders, expected %d", len(v.order), len(v.encs))
	}

	logs.Info("setup view", zap.Strings("encoders", v.order))
	return v, nil
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, d interface{}, opts ...Option) (err error) {
	if len(v.order) < 1 {
		return fmt.Errorf("no encoders configured")
	}

	fmt.Println(v.order)

	best, _ := accept.Negotiate(r.Header.Values("Accept"), v.order)
	if best < 0 {
		best = 0
	}

	enc, ok := v.encs[v.order[best]]
	if !ok {
		return fmt.Errorf("best matched MIME '%s' doesn't have a matching encoder", v.order[best])
	}

	var o Options
	for _, opt := range opts {
		opt(&o)
	}

	return enc.Encode(w, r, d, o)
}
