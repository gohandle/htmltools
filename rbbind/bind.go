package rbbind

import (
	"fmt"
	"mime"
	"net/http"
	"sort"

	"github.com/caarlos0/env/v6"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Conf configures the binder
type Conf struct {
	DecoderOrder []string `env:"RB_BIND_DECODER_ORDER" envSeparator:","`
}

// ParseConf parses the env
func ParseConf() (cfg Conf, err error) {
	return cfg, env.Parse(&cfg)
}

// Decoder can be implemented to decode http requests
type Decoder interface {
	Name() string
	Decode(r *http.Request, v interface{}, mt string, mtParams map[string]string) error
}

// Params configures dependencies binding dependencies
type Params struct {
	fx.In
	Decoders []Decoder `group:"rb.decoder"`
}

// Binder provides the ability to bind requests to values
type Binder struct {
	decs  map[string]Decoder
	order []string
	logs  *zap.Logger
}

// New inits the binder
func New(logs *zap.Logger, cfg Conf, p Params) (b *Binder, err error) {
	b = &Binder{
		decs:  make(map[string]Decoder, len(p.Decoders)),
		order: make([]string, 0, len(p.Decoders)),
	}
	for _, dec := range p.Decoders {
		if _, ok := b.decs[dec.Name()]; ok {
			return nil, fmt.Errorf("decoder with name '%s' already exists", dec.Name())
		}

		b.decs[dec.Name()] = dec
		b.order = append(b.order, dec.Name())
	}

	sort.Strings(b.order)
	if len(cfg.DecoderOrder) > 0 {
		b.order = cfg.DecoderOrder
	}

	if len(b.order) != len(b.decs) {
		return nil, fmt.Errorf("decoder order specified %d encoders, expected %d", len(b.order), len(b.decs))
	}

	logs.Info("setup binder", zap.Strings("decoders", b.order))
	return b, nil
}

// Bind will bind the request to all the provided values
func (b *Binder) Bind(r *http.Request, vs ...interface{}) (err error) {
	mt, mtParams, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	for _, v := range vs {
		for _, dn := range b.order {
			dec, ok := b.decs[dn]
			if !ok {
				return fmt.Errorf("decoder '%s' is in order configuration but not provided", dn)
			}

			err = dec.Decode(r, v, mt, mtParams)
			if err != nil {
				return fmt.Errorf("failed to decode: %w", err)
			}
		}
	}

	return nil
}
