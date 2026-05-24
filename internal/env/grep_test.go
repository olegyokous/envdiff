package env

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempGrepEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestGrep_MatchesKey(t *testing.T) {
	f := writeTempGrepEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=myapp\n")
	opts := DefaultGrepOptions()
	opts.Pattern = "^DB_"
	matches, err := Grep([]string{f}, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
}

func TestGrep_MatchesValue(t *testing.T) {
	f := writeTempGrepEnv(t, "HOST=localhost\nPORT=5432\n")
	opts := DefaultGrepOptions()
	opts.SearchKeys = false
	opts.Pattern = "local"
	matches, err := Grep([]string{f}, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 || matches[0].Key != "HOST" {
		t.Fatalf("unexpected matches: %+v", matches)
	}
}

func TestGrep_IgnoreCase(t *testing.T) {
	f := writeTempGrepEnv(t, "Secret=abc\nPUBLIC=xyz\n")
	opts := DefaultGrepOptions()
	opts.Pattern = "secret"
	opts.IgnoreCase = true
	matches, err := Grep([]string{f}, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
}

func TestGrep_EmptyPattern(t *testing.T) {
	f := writeTempGrepEnv(t, "KEY=val\n")
	_, err := Grep([]string{f}, GrepOptions{})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestGrep_InvalidPattern(t *testing.T) {
	f := writeTempGrepEnv(t, "KEY=val\n")
	opts := DefaultGrepOptions()
	opts.Pattern = "["
	_, err := Grep([]string{f}, opts)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestGrep_MissingFile(t *testing.T) {
	opts := DefaultGrepOptions()
	opts.Pattern = "KEY"
	_, err := Grep([]string{"/nonexistent/.env"}, opts)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestWriteGrepText_NoMatches(t *testing.T) {
	var buf bytes.Buffer
	WriteGrepText(&buf, nil)
	if !strings.Contains(buf.String(), "no matches") {
		t.Errorf("expected 'no matches' message, got: %s", buf.String())
	}
}

func TestWriteGrepJSON_ValidJSON(t *testing.T) {
	matches := []GrepMatch{{File: "a.env", Key: "FOO", Value: "bar"}}
	var buf bytes.Buffer
	if err := WriteGrepJSON(&buf, matches); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "FOO") {
		t.Errorf("expected key in JSON output: %s", buf.String())
	}
}

func TestWriteGrepSummary_ShowsCount(t *testing.T) {
	matches := []GrepMatch{
		{File: "a.env", Key: "A", Value: "1"},
		{File: "a.env", Key: "B", Value: "2"},
	}
	var buf bytes.Buffer
	WriteGrepSummary(&buf, matches)
	if !strings.Contains(buf.String(), "2 match") {
		t.Errorf("unexpected summary: %s", buf.String())
	}
}
