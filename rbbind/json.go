package rbbind

import (
	"encoding/json"
	"net/http"

	"go.uber.org/fx"
)

type JSON struct{}

func NewJSON() Decoder {
	return &JSON{}
}

func (d *JSON) Name() string { return "json" }
func (d *JSON) Decode(r *http.Request, v interface{}, mt string, _ map[string]string) (err error) {
	if mt != "application/json" {
		return nil
	}

	dec := json.NewDecoder(r.Body)
	return dec.Decode(v)
}

var JSONDecoder = fx.Annotated{
	Target: NewJSON,
	Group:  "rb.decoder",
}
