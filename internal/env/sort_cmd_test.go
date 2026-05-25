package env_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envdiff/internal/env"
)

func writeTempSortCmdEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempSortCmdEnv: %v", err)
	}
	return p
}

func TestRunSort_WritesFile(t *testing.T) {
	src := writeTempSortCmdEnv(t, "ZEBRA=1\nAPPLE=2\nMANGO=3\n")
	out := filepath.Join(t.TempDir(), "sorted.env")

	opts := env.DefaultSortOptions()
	opts.Source = src
	opts.Output = out

	if err := env.RunSort(opts); err != nil {
		t.Fatalf("RunSort: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "APPLE") {
		t.Errorf("expected first line APPLE, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA") {
		t.Errorf("expected last line ZEBRA, got %q", lines[2])
	}
}

func TestRunSort_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempSortCmdEnv(t, "ZEBRA=1\nAPPLE=2\n")
	out := filepath.Join(t.TempDir(), "should-not-exist.env")

	opts := env.DefaultSortOptions()
	opts.Source = src
	opts.Output = out
	opts.DryRun = true

	if err := env.RunSort(opts); err != nil {
		t.Fatalf("RunSort: %v", err)
	}

	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Error("expected output file not to be created in dry-run mode")
	}
}

func TestRunSort_MissingSource(t *testing.T) {
	opts := env.DefaultSortOptions()
	opts.Source = "/no/such/file.env"
	opts.Output = filepath.Join(t.TempDir(), "out.env")

	if err := env.RunSort(opts); err == nil {
		t.Error("expected error for missing source file")
	}
}

func TestRunSort_Descending(t *testing.T) {
	src := writeTempSortCmdEnv(t, "APPLE=1\nZEBRA=2\nMANGO=3\n")
	out := filepath.Join(t.TempDir(), "desc.env")

	opts := env.DefaultSortOptions()
	opts.Source = src
	opts.Output = out
	opts.Descending = true

	if err := env.RunSort(opts); err != nil {
		t.Fatalf("RunSort: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if !strings.HasPrefix(lines[0], "ZEBRA") {
		t.Errorf("expected first line ZEBRA in descending order, got %q", lines[0])
	}
}

func TestRunSort_InPlace(t *testing.T) {
	src := writeTempSortCmdEnv(t, "CHARLIE=3\nALPHA=1\nBRAVO=2\n")

	opts := env.DefaultSortOptions()
	opts.Source = src
	opts.Output = src // write back to source

	if err := env.RunSort(opts); err != nil {
		t.Fatalf("RunSort: %v", err)
	}

	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if !strings.HasPrefix(lines[0], "ALPHA") {
		t.Errorf("expected ALPHA first after in-place sort, got %q", lines[0])
	}
}
