package rblog

import (
	"errors"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingWare(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	RandRead = rnd.Read
	lc, obs := observer.New(zap.DebugLevel)

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/foo", nil)
	if nopl := L(r); nopl.Core().Enabled(zap.ErrorLevel) {
		t.Fatalf("should be nop")
	}

	h := NewLogging(zap.New(lc), CommonIDHeaders())(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		L(r).Debug("foo")
	}))

	h.ServeHTTP(w, r)
	if act := obs.FilterMessage("foo").All()[0].ContextMap()["request_id"]; act != "Uv38ByGCZU8WP18PmmIdcpVm" {
		t.Fatalf("got: %v", act)
	}

	t.Run("with id header in request", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/foo", nil)
		r.Header.Set("X-Amzn-Trace-Id", "foo")

		h.ServeHTTP(w, r)
		if act := obs.FilterMessage("foo").All()[1].ContextMap()["request_id"]; act != "foo" {
			t.Fatalf("got: %v", act)
		}
	})

	t.Run("failing rand read", func(t *testing.T) {
		RandRead = func(b []byte) (n int, err error) { return 0, errors.New("foo") }
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/foo", nil)

		h.ServeHTTP(w, r)
		if act := obs.FilterMessage("failed to read random bytes for request id middleware").All()[0].ContextMap()["error"]; act != "foo" {
			t.Fatalf("got: %v", act)
		}
	})
}
