package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/your-org/envdiff/internal/parser"
)

func writeSplitSrc(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "src*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestSplit_ByPrefix(t *testing.T) {
	src := writeSplitSrc(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=myapp\nAPP_ENV=prod\nOTHER=x\n")
	dir := t.TempDir()
	dbOut := filepath.Join(dir, "db.env")
	appOut := filepath.Join(dir, "app.env")
	restOut := filepath.Join(dir, "rest.env")

	counts, err := Split(src, SplitOptions{
		Prefixes:  map[string]string{"DB_": dbOut, "APP_": appOut},
		Remainder: restOut,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if counts[dbOut] != 2 {
		t.Errorf("db: want 2 keys, got %d", counts[dbOut])
	}
	if counts[appOut] != 2 {
		t.Errorf("app: want 2 keys, got %d", counts[appOut])
	}
	if counts[restOut] != 1 {
		t.Errorf("rest: want 1 key, got %d", counts[restOut])
	}
}

func TestSplit_StripPrefix(t *testing.T) {
	src := writeSplitSrc(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	dir := t.TempDir()
	dbOut := filepath.Join(dir, "db.env")

	_, err := Split(src, SplitOptions{
		Prefixes: map[string]string{"DB_": dbOut},
		Strip:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env, _ := parser.ParseFile(dbOut)
	if _, ok := env["HOST"]; !ok {
		t.Error("expected key HOST after stripping DB_ prefix")
	}
	if _, ok := env["DB_HOST"]; ok {
		t.Error("did not expect key DB_HOST when strip=true")
	}
}

func TestSplit_NoRemainderDiscards(t *testing.T) {
	src := writeSplitSrc(t, "DB_HOST=localhost\nUNMATCHED=foo\n")
	dir := t.TempDir()
	dbOut := filepath.Join(dir, "db.env")

	counts, err := Split(src, SplitOptions{
		Prefixes: map[string]string{"DB_": dbOut},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if counts[dbOut] != 1 {
		t.Errorf("want 1, got %d", counts[dbOut])
	}
	if len(counts) != 1 {
		t.Errorf("expected only 1 output file, got %d", len(counts))
	}
}

func TestSplit_MissingSource(t *testing.T) {
	_, err := Split("/nonexistent/file.env", SplitOptions{})
	if err == nil {
		t.Fatal("expected error for missing source file")
	}
}

func TestSplit_QuotedValuePreserved(t *testing.T) {
	src := writeSplitSrc(t, "APP_DESC=hello world\n")
	dir := t.TempDir()
	appOut := filepath.Join(dir, "app.env")

	_, err := Split(src, SplitOptions{
		Prefixes: map[string]string{"APP_": appOut},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	raw, _ := os.ReadFile(appOut)
	if !strings.Contains(string(raw), `"hello world"`) {
		t.Errorf("expected quoted value in output, got: %s", string(raw))
	}
}
