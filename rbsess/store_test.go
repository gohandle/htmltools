package rbsess

import (
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestParseConf(t *testing.T) {
	os.Setenv("RB_SESS_COOKIE_DOMAIN", "foo.example.com")
	defer os.Unsetenv("RB_SESS_COOKIE_DOMAIN")
	os.Setenv("RB_SESS_KEY_PAIRS", "xyz,def")
	defer os.Unsetenv("RB_SESS_KEY_PAIRS")

	cfg, err := ParseConf()
	if !reflect.DeepEqual(cfg, Conf{
		CookieDomain:        "foo.example.com",
		KeyPairs:            []string{"xyz", "def"},
		CookieMaxAgeSeconds: 2592000,
		CookiePath:          "/",
	}) {
		t.Fatalf("got: %+v err: %v", cfg, err)
	}
}

func TestStoreInit(t *testing.T) {
	logs, _ := zap.NewDevelopment()

	s, err := NewCookieStore(logs, Conf{KeyPairs: []string{"7Ijc+Qj3h4GLxAFd3q1dUQ==", "7Ij5+Qj3h4GLxAFd3q1dUQ=="}})
	if err != nil || s == nil {
		t.Fatalf("got: %v %v", err, s)
	}

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	sess, err := s.Get(r, "foo")
	if err != nil {
		t.Fatalf("failed to get: %v", err)
	}

	if err = s.Save(r, w, sess); err != nil {
		t.Fatalf("got: %v", err)
	}

	if !strings.HasPrefix(w.Header().Get("Set-Cookie"), "foo") {
		t.Fatalf("got: %v", w.Header().Get("Set-Cookie"))
	}
}
