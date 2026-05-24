package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/parser"
)

func writeTempSetEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempSetEnv: %v", err)
	}
	return p
}

func TestSet_UpdatesExistingKey(t *testing.T) {
	p := writeTempSetEnv(t, "FOO=bar\nBAZ=qux\n")
	res, err := Set(SetOptions{File: p, Key: "FOO", Value: "newval"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Created {
		t.Error("expected Created=false for existing key")
	}
	if res.Prev != "bar" {
		t.Errorf("expected prev=bar, got %q", res.Prev)
	}
	env, _ := parser.ParseFile(p)
	if env["FOO"] != "newval" {
		t.Errorf("expected FOO=newval on disk, got %q", env["FOO"])
	}
}

func TestSet_CreatesNewKey(t *testing.T) {
	p := writeTempSetEnv(t, "FOO=bar\n")
	res, err := Set(SetOptions{File: p, Key: "NEW", Value: "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Created {
		t.Error("expected Created=true for new key")
	}
	env, _ := parser.ParseFile(p)
	if env["NEW"] != "hello" {
		t.Errorf("expected NEW=hello on disk, got %q", env["NEW"])
	}
}

func TestSet_DryRunDoesNotWrite(t *testing.T) {
	p := writeTempSetEnv(t, "FOO=bar\n")
	_, err := Set(SetOptions{File: p, Key: "FOO", Value: "changed", DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env, _ := parser.ParseFile(p)
	if env["FOO"] != "bar" {
		t.Errorf("dry-run should not modify file, got %q", env["FOO"])
	}
}

func TestSet_MissingFileWithoutCreate(t *testing.T) {
	_, err := Set(SetOptions{File: "/nonexistent/.env", Key: "X", Value: "1"})
	if err == nil {
		t.Error("expected error for missing file without Create flag")
	}
}

func TestSet_MissingFileWithCreate(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "new.env")
	res, err := Set(SetOptions{File: p, Key: "X", Value: "1", Create: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Created {
		t.Error("expected Created=true")
	}
	env, _ := parser.ParseFile(p)
	if env["X"] != "1" {
		t.Errorf("expected X=1, got %q", env["X"])
	}
}

func TestSet_EmptyKeyReturnsError(t *testing.T) {
	p := writeTempSetEnv(t, "")
	_, err := Set(SetOptions{File: p, Key: "", Value: "v"})
	if err == nil {
		t.Error("expected error for empty key")
	}
}
