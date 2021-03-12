package rblog

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingWare(t *testing.T) {
	lc, obs := observer.New(zap.DebugLevel)

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/foo", nil)
	if nopl := L(r); nopl.Core().Enabled(zap.ErrorLevel) {
		t.Fatalf("should be nop")
	}

	NewLogging(zap.New(lc))(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		L(r).Debug("foo")
	})).ServeHTTP(w, r)

	if act := obs.FilterMessage("foo").All()[0].ContextMap()["request_uri"]; act != "/foo" {
		t.Fatalf("got: %v", act)
	}
}
