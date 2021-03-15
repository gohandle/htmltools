package rbbind

import (
	"encoding/json"
	"net/http"

	"go.uber.org/fx"
)

// JSON decoder
type JSON struct{}

// NewJSON creates the json request decoder
func NewJSON() Decoder {
	return &JSON{}
}

// Name returns the name of the json decoder
func (d *JSON) Name() string { return "json" }

// Decode performs the request decoding
func (d *JSON) Decode(r *http.Request, v interface{}, mt string, _ map[string]string) (err error) {
	if mt != "application/json" {
		return nil
	}

	dec := json.NewDecoder(r.Body)
	return dec.Decode(v)
}

// JSONDecoder annotates the json decoder constructor as the correct group
var JSONDecoder = fx.Annotated{
	Target: NewJSON,
	Group:  "rb.decoder",
}
