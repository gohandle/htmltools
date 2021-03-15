package i18n

import (
	"embed"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestParseConf(t *testing.T) {
	os.Setenv("RB_I18N_DIR", "foo")
	defer os.Unsetenv("RB_I18N_DIR")
	os.Setenv("RB_I18N_FILE_NAMES", "foo.en-US.json:bar.en-US.toml")
	defer os.Unsetenv("RB_I18N_FILE_NAMES")

	cfg, err := ParseConf()
	if !reflect.DeepEqual(cfg, Conf{
		Dir:             "foo",
		DefaultLanguage: "en-US",
		FileNames:       []string{"foo.en-US.json", "bar.en-US.toml"},
	}) {
		t.Fatalf("got: %+v err: %v", cfg, err)
	}
}

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
	var b *i18n.Bundle
	fxtest.New(t, fx.Populate(&b),
		fx.Supply(testFs, Conf{DefaultLanguage: "en-US", FileNames: []string{"test.en-US.json"}}),
		fx.Provide(FromEmbed, New, zap.NewDevelopment)).RequireStart().RequireStop()

	s, err := i18n.NewLocalizer(b, "en").Localize(&i18n.LocalizeConfig{MessageID: "foo"})
	if err != nil || s != "rab" {
		t.Fatalf("got: %v (%v)", err, s)
	}
}
