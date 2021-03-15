package i18n

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
)

//go:embed test.en-US.json
var testFs embed.FS

func TestI18NInitWithoutMessages(t *testing.T) {
	logs, _ := zap.NewDevelopment()

	t.Run("invalid default language", func(t *testing.T) {
		_, err := New(logs, Conf{DefaultLanguage: "xy-df"}, Params{})
		if !strings.Contains(err.Error(), "invalid default language") {
			t.Fatalf("got: %v", err)
		}
	})

	b, err := New(logs, Conf{DefaultLanguage: "en-US"}, Params{})
	if err != nil || b == nil {
		t.Fatalf("got: %v %v", err, b)
	}
}

func TestI18NInitWithMessagesFromOsDir(t *testing.T) {
	logs, _ := zap.NewDevelopment()
	dir, _ := os.MkdirTemp("", "")
	os.WriteFile(filepath.Join(dir, "msgs.en-US.json"), []byte(`{"foo": "bar"}`), 0777)
	cfg := Conf{DefaultLanguage: "en-US", Dir: dir, FileNames: []string{"msgs.en-US.json"}}

	b, err := New(logs, cfg, Params{Files: FromDir(logs, cfg)})
	if err != nil || b == nil {
		t.Fatalf("got: %v %v", err, b)
	}

	s, err := i18n.NewLocalizer(b, "en").Localize(&i18n.LocalizeConfig{MessageID: "foo"})
	if err != nil || s != "bar" {
		t.Fatalf("got: %v (%v)", err, s)
	}
}

func TestI18NInitWithMessagesFromEmbed(t *testing.T) {
	logs, _ := zap.NewDevelopment()

	cfg := Conf{DefaultLanguage: "en-US", FileNames: []string{"test.en-US.json"}}

	b, err := New(logs, cfg, Params{Files: FromEmbed(testFs)})
	if err != nil || b == nil {
		t.Fatalf("got: %v %v", err, b)
	}

	s, err := i18n.NewLocalizer(b, "en").Localize(&i18n.LocalizeConfig{MessageID: "foo"})
	if err != nil || s != "rab" {
		t.Fatalf("got: %v (%v)", err, s)
	}
}
