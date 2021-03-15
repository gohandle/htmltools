package rbbind

import (
	"fmt"
	"net/http"

	"github.com/go-playground/form/v4"
	"go.uber.org/fx"
)

// Form decoder supports decoding of url encoded values
type Form struct{ dec *form.Decoder }

// NewForm inits the form decoder for binding
func NewForm() Decoder {
	return &Form{dec: form.NewDecoder()}
}

// Name returns the name of the decoder
func (d *Form) Name() string { return "form" }

// Decode performs the actual decoder
func (d *Form) Decode(r *http.Request, v interface{}, _ string, _ map[string]string) (err error) {
	if err = r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form: %w", err)
	}

	return d.dec.Decode(v, r.Form)
}

// FormDecoder provides an form decoder annotated with the correct group
var FormDecoder = fx.Annotated{
	Target: NewForm,
	Group:  "rb.decoder",
}
