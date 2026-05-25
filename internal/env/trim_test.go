package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempTrimEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestTrim_RemovesEmptyValues(t *testing.T) {
	src := writeTempTrimEnv(t, "KEEP=hello\nREMOVE=\nALSO_KEEP=world\n")
	out := filepath.Join(t.TempDir(), "out.env")

	result, err := Trim(TrimOptions{Source: src, Output: out, EmptyOnly: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	removed := countRemoved(result.Records)
	if removed != 1 {
		t.Errorf("expected 1 removed, got %d", removed)
	}

	data, _ := os.ReadFile(out)
	if contains(string(data), "REMOVE") {
		t.Error("output should not contain REMOVE key")
	}
}

func TestTrim_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempTrimEnv(t, "A=\nB=value\n")
	out := filepath.Join(t.TempDir(), "out.env")

	_, err := Trim(TrimOptions{Source: src, Output: out, DryRun: true, EmptyOnly: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Error("dry run should not create output file")
	}
}

func TestTrim_AllKeysWhenNotEmptyOnly(t *testing.T) {
	src := writeTempTrimEnv(t, "A=one\nB=two\n")
	out := filepath.Join(t.TempDir(), "out.env")

	result, err := Trim(TrimOptions{Source: src, Output: out, EmptyOnly: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if countRemoved(result.Records) != 2 {
		t.Errorf("expected all 2 keys removed")
	}

	data, _ := os.ReadFile(out)
	if len(data) != 0 {
		t.Errorf("expected empty output file, got %q", string(data))
	}
}

func TestTrim_MissingSource(t *testing.T) {
	_, err := Trim(TrimOptions{Source: "/no/such/file.env"})
	if err == nil {
		t.Error("expected error for missing source")
	}
}

func TestTrim_WhitespaceOnlyValueRemoved(t *testing.T) {
	src := writeTempTrimEnv(t, "A=   \nB=real\n")
	out := filepath.Join(t.TempDir(), "out.env")

	result, err := Trim(TrimOptions{Source: src, Output: out, EmptyOnly: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if countRemoved(result.Records) != 1 {
		t.Errorf("expected 1 removed (whitespace-only), got %d", countRemoved(result.Records))
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
