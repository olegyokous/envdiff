package env

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempDedupeEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestDedupe_NoDuplicates(t *testing.T) {
	src := writeTempDedupeEnv(t, "FOO=1\nBAR=2\n")
	res, err := Dedupe(src, DedupeOptions{KeepFirst: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Duplicates) != 0 {
		t.Errorf("expected no duplicates, got %v", res.Duplicates)
	}
}

func TestDedupe_KeepFirst(t *testing.T) {
	src := writeTempDedupeEnv(t, "FOO=first\nBAR=1\nFOO=second\n")
	res, err := Dedupe(src, DedupeOptions{KeepFirst: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Duplicates) != 1 || res.Duplicates[0] != "FOO" {
		t.Errorf("expected [FOO] duplicates, got %v", res.Duplicates)
	}
	data, _ := os.ReadFile(src)
	if !strings.Contains(string(data), "FOO=first") {
		t.Error("expected first value to be kept")
	}
	if strings.Contains(string(data), "FOO=second") {
		t.Error("second value should have been removed")
	}
}

func TestDedupe_KeepLast(t *testing.T) {
	src := writeTempDedupeEnv(t, "FOO=first\nBAR=1\nFOO=second\n")
	res, err := Dedupe(src, DedupeOptions{KeepFirst: false})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(src)
	if strings.Contains(string(data), "FOO=first") {
		t.Error("first value should have been removed")
	}
	if !strings.Contains(string(data), "FOO=second") {
		t.Error("expected last value to be kept")
	}
	_ = res
}

func TestDedupe_DryRunDoesNotWrite(t *testing.T) {
	original := "FOO=a\nFOO=b\n"
	src := writeTempDedupeEnv(t, original)
	_, err := Dedupe(src, DedupeOptions{KeepFirst: true, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(src)
	if string(data) != original {
		t.Error("dry-run should not modify the file")
	}
}

func TestDedupe_MissingFile(t *testing.T) {
	_, err := Dedupe("/no/such/file.env", DedupeOptions{})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestWriteDedupeText_NoDuplicates(t *testing.T) {
	r := &DedupeResult{Source: "test.env", Duplicates: nil}
	var buf bytes.Buffer
	WriteDedupeText(&buf, r)
	if !strings.Contains(buf.String(), "no duplicate") {
		t.Error("expected 'no duplicate' message")
	}
}

func TestWriteDedupeJSON_ValidJSON(t *testing.T) {
	r := &DedupeResult{Source: "a.env", Duplicates: []string{"FOO", "BAR"}}
	var buf bytes.Buffer
	if err := WriteDedupeJSON(&buf, r); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"duplicate_count": 2`) {
		t.Error("expected duplicate_count in JSON output")
	}
}
