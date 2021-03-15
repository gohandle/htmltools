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

type Conf struct {
	DecoderOrder []string `env:"RB_BIND_DECODER_ORDER" envSeparator:","`
}

// ParseConf parses the env
func ParseConf() (cfg Conf, err error) {
	return cfg, env.Parse(&cfg)
}

type Decoder interface {
	Name() string
	Decode(r *http.Request, v interface{}, mt string, mtParams map[string]string) error
}

type Binder struct {
	decs  map[string]Decoder
	order []string
}

type Params struct {
	fx.In
	Decoders []Decoder `group:"rb.decoder"`
}

func New(logs *zap.Logger, cfg Conf, p Params) (b *Binder, err error) {
	b = &Binder{
		decs:  make(map[string]Decoder, len(p.Decoders)),
		order: make([]string, 0, len(p.Decoders)),
	}
	for _, dec := range p.Decoders {
		b.decs[dec.Name()] = dec
		b.order = append(b.order, dec.Name())
	}

	sort.Strings(b.order)
	if len(cfg.DecoderOrder) > 0 {
		b.order = cfg.DecoderOrder
	}

	return b, nil
}

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
