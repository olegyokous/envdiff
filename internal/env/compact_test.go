package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempCompactEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp compact env: %v", err)
	}
	return p
}

func TestCompact_RemovesCommentsAndBlanks(t *testing.T) {
	src := writeTempCompactEnv(t, "# comment\nFOO=bar\n\nBAZ=qux\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	result, err := Compact(src, dst, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RemovedCount != 2 {
		t.Errorf("expected 2 removed, got %d", result.RemovedCount)
	}
	if result.KeptCount != 2 {
		t.Errorf("expected 2 kept, got %d", result.KeptCount)
	}
}

func TestCompact_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempCompactEnv(t, "# comment\nFOO=bar\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	_, err := Compact(src, dst, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		t.Error("expected output file not to exist in dry-run mode")
	}
}

func TestCompact_WritesOnlyKeyValues(t *testing.T) {
	src := writeTempCompactEnv(t, "# header\nFOO=bar\n\nBAZ=qux\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	_, err := Compact(src, dst, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if strings.HasPrefix(line, "#") || line == "" {
			t.Errorf("unexpected line in output: %q", line)
		}
	}
}

func TestCompact_MissingSource(t *testing.T) {
	_, err := Compact("/no/such/file.env", "/tmp/out.env", false)
	if err == nil {
		t.Error("expected error for missing source")
	}
}

func TestCompact_NoCommentsOrBlanks(t *testing.T) {
	src := writeTempCompactEnv(t, "FOO=bar\nBAZ=qux\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	result, err := Compact(src, dst, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RemovedCount != 0 {
		t.Errorf("expected 0 removed, got %d", result.RemovedCount)
	}
}
