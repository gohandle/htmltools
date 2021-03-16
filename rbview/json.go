package rbview

import (
	"encoding/json"
	"net/http"

	"go.uber.org/fx"
)

type JSON struct{}

func NewJSON() Encoder {
	return &JSON{}
}

func (e *JSON) MIME() string { return "application/json" }
func (e *JSON) Encode(w http.ResponseWriter, r *http.Request, v interface{}, o Options) error {
	return json.NewEncoder(w).Encode(v)
}

var JSONEncoder = fx.Annotated{
	Target: NewJSON,
	Group:  "rb.encoder",
}
