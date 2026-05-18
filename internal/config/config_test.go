package config

import (
	"testing"
)

func TestParse_MinimalValid(t *testing.T) {
	opts, err := Parse([]string{"a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(opts.Files))
	}
	if opts.Format != "text" {
		t.Errorf("expected default format \"text\", got %q", opts.Format)
	}
}

func TestParse_TooFewFiles(t *testing.T) {
	_, err := Parse([]string{"only.env"})
	if err == nil {
		t.Fatal("expected error for fewer than 2 files")
	}
}

func TestParse_NoFiles(t *testing.T) {
	_, err := Parse([]string{})
	if err == nil {
		t.Fatal("expected error when no files provided")
	}
}

func TestParse_JSONFormat(t *testing.T) {
	opts, err := Parse([]string{"-format", "json", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Format != "json" {
		t.Errorf("expected format \"json\", got %q", opts.Format)
	}
}

func TestParse_InvalidFormat(t *testing.T) {
	_, err := Parse([]string{"-format", "yaml", "a.env", "b.env"})
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestParse_StatusFilter(t *testing.T) {
	opts, err := Parse([]string{"-status", "missing", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.StatusFilter != "missing" {
		t.Errorf("expected status filter \"missing\", got %q", opts.StatusFilter)
	}
}

func TestParse_InvalidStatusFilter(t *testing.T) {
	_, err := Parse([]string{"-status", "unknown", "a.env", "b.env"})
	if err == nil {
		t.Fatal("expected error for invalid status filter")
	}
}

func TestParse_Flags(t *testing.T) {
	opts, err := Parse([]string{
		"-no-color",
		"-exit-code",
		"-prefix", "APP_",
		"-pattern", "^DB_",
		"a.env", "b.env", "c.env",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.NoColor {
		t.Error("expected NoColor to be true")
	}
	if !opts.ExitCode {
		t.Error("expected ExitCode to be true")
	}
	if opts.KeyPrefix != "APP_" {
		t.Errorf("expected prefix \"APP_\", got %q", opts.KeyPrefix)
	}
	if opts.KeyPattern != "^DB_" {
		t.Errorf("expected pattern \"^DB_\", got %q", opts.KeyPattern)
	}
	if len(opts.Files) != 3 {
		t.Errorf("expected 3 files, got %d", len(opts.Files))
	}
}
