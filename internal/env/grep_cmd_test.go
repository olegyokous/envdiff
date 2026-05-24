package env_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envdiff/internal/env"
)

func writeGrepCmdEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeGrepCmdEnv: %v", err)
	}
	return p
}

func TestRunGrep_MatchesKeyOutput(t *testing.T) {
	p := writeGrepCmdEnv(t, "DATABASE_URL=postgres://localhost\nSECRET_KEY=abc123\nAPP_PORT=8080\n")

	var buf bytes.Buffer
	opts := env.DefaultGrepCmdOptions()
	opts.Pattern = "DATABASE"
	opts.Files = []string{p}
	opts.Out = &buf

	if err := env.RunGrep(opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DATABASE_URL") {
		t.Errorf("expected DATABASE_URL in output, got:\n%s", out)
	}
	if strings.Contains(out, "SECRET_KEY") {
		t.Errorf("did not expect SECRET_KEY in output, got:\n%s", out)
	}
}

func TestRunGrep_JSONFormat(t *testing.T) {
	p := writeGrepCmdEnv(t, "API_KEY=secret\nAPI_SECRET=topsecret\nHOST=localhost\n")

	var buf bytes.Buffer
	opts := env.DefaultGrepCmdOptions()
	opts.Pattern = "API"
	opts.Files = []string{p}
	opts.Format = "json"
	opts.Out = &buf

	if err := env.RunGrep(opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in JSON output, got:\n%s", out)
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "[") {
		t.Errorf("expected JSON array output, got:\n%s", out)
	}
}

func TestRunGrep_SummaryFormat(t *testing.T) {
	p := writeGrepCmdEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nDB_NAME=mydb\nAPP_ENV=production\n")

	var buf bytes.Buffer
	opts := env.DefaultGrepCmdOptions()
	opts.Pattern = "DB_"
	opts.Files = []string{p}
	opts.Format = "summary"
	opts.Out = &buf

	if err := env.RunGrep(opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "3") {
		t.Errorf("expected match count 3 in summary, got:\n%s", out)
	}
}

func TestRunGrep_MissingFile(t *testing.T) {
	var buf bytes.Buffer
	opts := env.DefaultGrepCmdOptions()
	opts.Pattern = "KEY"
	opts.Files = []string{"/nonexistent/.env"}
	opts.Out = &buf

	if err := env.RunGrep(opts); err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestRunGrep_InvalidPattern(t *testing.T) {
	p := writeGrepCmdEnv(t, "KEY=value\n")

	var buf bytes.Buffer
	opts := env.DefaultGrepCmdOptions()
	opts.Pattern = "[invalid"
	opts.Files = []string{p}
	opts.Out = &buf

	if err := env.RunGrep(opts); err == nil {
		t.Error("expected error for invalid regex pattern, got nil")
	}
}

func TestRunGrep_MultipleFiles(t *testing.T) {
	p1 := writeGrepCmdEnv(t, "SHARED_KEY=value1\nONLY_IN_A=yes\n")
	p2 := writeGrepCmdEnv(t, "SHARED_KEY=value2\nONLY_IN_B=yes\n")

	var buf bytes.Buffer
	opts := env.DefaultGrepCmdOptions()
	opts.Pattern = "SHARED"
	opts.Files = []string{p1, p2}
	opts.Out = &buf

	if err := env.RunGrep(opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	// Both files should produce a match
	count := strings.Count(out, "SHARED_KEY")
	if count < 2 {
		t.Errorf("expected SHARED_KEY to appear at least twice (once per file), got %d occurrences in:\n%s", count, out)
	}
}
