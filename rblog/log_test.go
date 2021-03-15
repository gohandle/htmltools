package rblog

import (
	"os"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestConf(t *testing.T) {
	os.Setenv("RB_LOG_DEV_MODE", "true")
	defer os.Unsetenv("RB_LOG_DEV_MODE")
	cfg, err := ParseConf()
	if err != nil {
		t.Fatalf("got: %v", err)
	}

	if !cfg.DevMode {
		t.Fatalf("got: %v", cfg.DevMode)
	}
}

func TestNewStd(t *testing.T) {
	l := NewStd()
	if l == nil {
		t.Fatalf("got: %v", l)
	}
}

func TestNew(t *testing.T) {
	t.Run("prod", func(t *testing.T) {
		l1, err := New(Conf{})
		if err != nil {
			t.Fatalf("got: %v", err)
		}

		if act := l1.Core().Enabled(zap.DebugLevel); act {
			t.Fatalf("got:  %v", act)
		}
	})

	t.Run("dev", func(t *testing.T) {
		l2, err := New(Conf{DevMode: true})
		if err != nil {
			t.Fatalf("got: %v", err)
		}

		if act := l2.Core().Enabled(zap.DebugLevel); !act {
			t.Fatalf("got:  %v", act)
		}
	})
}

func TestFxIntegration(t *testing.T) {
	var mw LoggingWare

	fxtest.New(t, fx.Provide(
		zap.NewDevelopment,
		NewLogging,
		CommonIDHeaders,
	), fx.Populate(&mw)).RequireStart().RequireStop()
	if mw == nil {
		t.Fatalf("got: %v", mw)
	}
}
