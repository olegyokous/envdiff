package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return path
}

func TestParseFile_Basic(t *testing.T) {
	path := writeTemp(t, "KEY=value\nFOO=bar\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", env["KEY"])
	}
	if env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", env["FOO"])
	}
}

func TestParseFile_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTemp(t, "# comment\n\nKEY=value\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 key, got %d", len(env))
	}
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTemp(t, `KEY="hello world"
FOO='single'`)
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", env["KEY"])
	}
	if env["FOO"] != "single" {
		t.Errorf("expected 'single', got %q", env["FOO"])
	}
}

func TestParseFile_MissingFile(t *testing.T) {
	_, err := ParseFile("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestParseFile_EmptyKey(t *testing.T) {
	path := writeTemp(t, "=value\n")
	_, err := ParseFile(path)
	if err == nil {
		t.Error("expected error for empty key, got nil")
	}
}

func TestStripQuotes(t *testing.T) {
	cases := []struct{ in, want string }{
		{`"quoted"`, "quoted"},
		{`'single'`, "single"},
		{`noquote`, "noquote"},
		{`"`, `"`},
	}
	for _, c := range cases {
		got := stripQuotes(c.in)
		if got != c.want {
			t.Errorf("stripQuotes(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
