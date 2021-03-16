package rbview

import "net/http"

type Encoder interface {
	Encode(w http.ResponseWriter, r *http.Request) error
}

type View struct {
	encs []Encoder
}
