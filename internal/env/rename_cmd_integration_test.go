package env_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/env"
	"github.com/yourorg/envdiff/internal/parser"
)

func writeIntegrationEnv(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRenameKey_VerboseOutput(t *testing.T) {
	src := writeIntegrationEnv(t, ".env", "SECRET_KEY=abc123\nDEBUG=true\n")
	opts := env.DefaultRenameOptions()
	opts.Source = src
	opts.OldKey = "SECRET_KEY"
	opts.NewKey = "APP_SECRET"
	opts.Verbose = true

	var buf bytes.Buffer
	if err := env.RenameKey(opts, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "SECRET_KEY -> APP_SECRET") {
		t.Errorf("expected verbose output, got: %s", buf.String())
	}
}

func TestRenameKey_PreservesOtherKeys(t *testing.T) {
	src := writeIntegrationEnv(t, ".env", "A=1\nB=2\nC=3\n")
	opts := env.DefaultRenameOptions()
	opts.Source = src
	opts.OldKey = "B"
	opts.NewKey = "BRAVO"

	if err := env.RenameKey(opts, &bytes.Buffer{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := parser.ParseFile(src)
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"A", "BRAVO", "C"} {
		if _, ok := result[k]; !ok {
			t.Errorf("expected key %q to be present", k)
		}
	}
	if _, ok := result["B"]; ok {
		t.Error("expected B to be removed after rename")
	}
}

func TestRenameKey_MissingSourceFile(t *testing.T) {
	opts := env.DefaultRenameOptions()
	opts.Source = "/nonexistent/.env"
	opts.OldKey = "FOO"
	opts.NewKey = "BAR"

	err := env.RenameKey(opts, &bytes.Buffer{})
	if err == nil {
		t.Error("expected error for missing source file")
	}
}
