package config_test

import (
	"flag"
	"testing"

	"github.com/your-org/envdiff/internal/config"
	"github.com/your-org/envdiff/internal/redact"
)

func newFS() *flag.FlagSet {
	return flag.NewFlagSet("test", flag.ContinueOnError)
}

func TestParse_MinimalValid(t *testing.T) {
	cfg, err := config.Parse(newFS(), []string{"a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(cfg.Files))
	}
}

func TestParse_TooFewFiles(t *testing.T) {
	_, err := config.Parse(newFS(), []string{"only.env"})
	if err == nil {
		t.Error("expected error for fewer than 2 files")
	}
}

func TestParse_NoFiles(t *testing.T) {
	_, err := config.Parse(newFS(), []string{})
	if err == nil {
		t.Error("expected error for no files")
	}
}

func TestParse_JSONFormat(t *testing.T) {
	cfg, err := config.Parse(newFS(), []string{"-format", "json", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format.String() != "json" {
		t.Errorf("expected json format, got %s", cfg.Format)
	}
}

func TestParse_InvalidFormat(t *testing.T) {
	_, err := config.Parse(newFS(), []string{"-format", "xml", "a.env", "b.env"})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestParse_RedactDefaultsIncluded(t *testing.T) {
	cfg, err := config.Parse(newFS(), []string{"a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	l := cfg.RedactList
	for _, p := range redact.DefaultSensitivePatterns {
		if !l.IsSensitive(p + "_FIELD") {
			t.Errorf("expected default pattern %q to be sensitive", p)
		}
	}
}

func TestParse_RedactCustomKeys(t *testing.T) {
	cfg, err := config.Parse(newFS(), []string{"-redact", "INTERNAL,LEGACY", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.RedactList.IsSensitive("INTERNAL_HOST") {
		t.Error("expected INTERNAL_HOST to be sensitive after custom redact flag")
	}
	if !cfg.RedactList.IsSensitive("LEGACY_TOKEN") {
		t.Error("expected LEGACY_TOKEN to be sensitive after custom redact flag")
	}
}

func TestParse_Flags(t *testing.T) {
	cfg, err := config.Parse(newFS(), []string{
		"-status", "missing",
		"-key-prefix", "APP_",
		"-no-summary",
		"a.env", "b.env",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.StatusFilter != "missing" {
		t.Errorf("expected status filter 'missing', got %q", cfg.StatusFilter)
	}
	if cfg.KeyPrefix != "APP_" {
		t.Errorf("expected key prefix 'APP_', got %q", cfg.KeyPrefix)
	}
	if !cfg.NoSummary {
		t.Error("expected NoSummary to be true")
	}
}

func TestParse_ThreeOrMoreFiles(t *testing.T) {
	cfg, err := config.Parse(newFS(), []string{"a.env", "b.env", "c.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Files) != 3 {
		t.Errorf("expected 3 files, got %d", len(cfg.Files))
	}
}
