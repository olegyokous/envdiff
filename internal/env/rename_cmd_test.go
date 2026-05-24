package env

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/parser"
)

func writeTempRenameEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRenameKey_BasicRename(t *testing.T) {
	src := writeTempRenameEnv(t, "OLD_KEY=hello\nOTHER=world\n")
	opts := DefaultRenameOptions()
	opts.Source = src
	opts.OldKey = "OLD_KEY"
	opts.NewKey = "NEW_KEY"

	var buf bytes.Buffer
	if err := RenameKey(opts, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	env, err := parser.ParseFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := env["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if v, ok := env["NEW_KEY"]; !ok || v != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", v)
	}
}

func TestRenameKey_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempRenameEnv(t, "FOO=bar\n")
	opts := DefaultRenameOptions()
	opts.Source = src
	opts.OldKey = "FOO"
	opts.NewKey = "BAZ"
	opts.DryRun = true

	var buf bytes.Buffer
	if err := RenameKey(opts, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	env, _ := parser.ParseFile(src)
	if _, ok := env["FOO"]; !ok {
		t.Error("dry-run should not modify source file")
	}
	if !strings.Contains(buf.String(), "FOO -> BAZ") {
		t.Errorf("expected dry-run output, got: %s", buf.String())
	}
}

func TestRenameKey_MissingOldKey(t *testing.T) {
	src := writeTempRenameEnv(t, "FOO=bar\n")
	opts := DefaultRenameOptions()
	opts.Source = src
	opts.OldKey = "MISSING"
	opts.NewKey = "NEW"

	err := RenameKey(opts, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error, got: %v", err)
	}
}

func TestRenameKey_ConflictingNewKey(t *testing.T) {
	src := writeTempRenameEnv(t, "FOO=bar\nBAZ=qux\n")
	opts := DefaultRenameOptions()
	opts.Source = src
	opts.OldKey = "FOO"
	opts.NewKey = "BAZ"

	err := RenameKey(opts, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' error, got: %v", err)
	}
}

func TestRenameKey_WritesToDest(t *testing.T) {
	src := writeTempRenameEnv(t, "ALPHA=1\nBETA=2\n")
	dest := filepath.Join(t.TempDir(), "out.env")

	opts := DefaultRenameOptions()
	opts.Source = src
	opts.Dest = dest
	opts.OldKey = "ALPHA"
	opts.NewKey = "GAMMA"

	if err := RenameKey(opts, &bytes.Buffer{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	env, err := parser.ParseFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := env["GAMMA"]; !ok || v != "1" {
		t.Errorf("expected GAMMA=1 in dest, got %q", v)
	}
	// source should be unchanged
	srcEnv, _ := parser.ParseFile(src)
	if _, ok := srcEnv["ALPHA"]; !ok {
		t.Error("source file should be unchanged when dest is different")
	}
}

func TestRenameKey_EmptyKeyErrors(t *testing.T) {
	opts := DefaultRenameOptions()
	opts.Source = "irrelevant.env"
	opts.OldKey = ""
	opts.NewKey = "NEW"

	err := RenameKey(opts, &bytes.Buffer{})
	if err == nil {
		t.Error("expected error for empty old-key")
	}
}
