package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempMaskEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestMask_MasksSpecifiedKeys(t *testing.T) {
	src := writeTempMaskEnv(t, "SECRET=hunter2\nPUBLIC=hello\n")
	out := filepath.Join(t.TempDir(), "out.env")
	opts := MaskOptions{Keys: []string{"SECRET"}, Output: out}
	res, err := Mask(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	maskedCount := 0
	for _, r := range res.Records {
		if r.Masked {
			maskedCount++
			if r.Key != "SECRET" {
				t.Errorf("unexpected masked key: %s", r.Key)
			}
		}
	}
	if maskedCount != 1 {
		t.Errorf("expected 1 masked key, got %d", maskedCount)
	}
}

func TestMask_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempMaskEnv(t, "TOKEN=abc123\n")
	out := filepath.Join(t.TempDir(), "should_not_exist.env")
	opts := MaskOptions{Keys: []string{"TOKEN"}, DryRun: true, Output: out}
	_, err := Mask(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Error("expected output file to not be created during dry run")
	}
}

func TestMask_FixedLength(t *testing.T) {
	src := writeTempMaskEnv(t, "PWD=supersecret\n")
	out := filepath.Join(t.TempDir(), "out.env")
	opts := MaskOptions{Keys: []string{"PWD"}, MaskLen: 6, Output: out}
	_, err := Mask(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !containsStr(string(data), "PWD=******") {
		t.Errorf("expected fixed-length mask, got: %s", string(data))
	}
}

func TestMask_MissingSource(t *testing.T) {
	_, err := Mask("/nonexistent/.env", MaskOptions{Keys: []string{"X"}})
	if err == nil {
		t.Error("expected error for missing source file")
	}
}

func TestMask_NoKeysMatchedLeavesAll(t *testing.T) {
	src := writeTempMaskEnv(t, "A=1\nB=2\n")
	out := filepath.Join(t.TempDir(), "out.env")
	opts := MaskOptions{Keys: []string{"NOPE"}, Output: out}
	res, err := Mask(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range res.Records {
		if r.Masked {
			t.Errorf("expected no masked keys, but %s was masked", r.Key)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsSubstring(s, sub))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
