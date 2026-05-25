package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempRotateEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempRotateEnv: %v", err)
	}
	return p
}

func TestRotate_RotatesKey(t *testing.T) {
	src := writeTempRotateEnv(t, "API_KEY=old_secret\nDEBUG=true\n")
	dst := src

	result, err := Rotate(src, dst, RotateOptions{
		Keys:   []string{"API_KEY"},
		Values: map[string]string{"API_KEY": "new_secret"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(result.Records))
	}
	rec := result.Records[0]
	if rec.Key != "API_KEY" || rec.OldValue != "old_secret" || rec.NewValue != "new_secret" || !rec.Rotated {
		t.Errorf("unexpected record: %+v", rec)
	}
}

func TestRotate_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempRotateEnv(t, "TOKEN=abc\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	_, err := Rotate(src, dst, RotateOptions{
		Keys:   []string{"TOKEN"},
		Values: map[string]string{"TOKEN": "xyz"},
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		t.Error("expected dst file to not exist in dry-run mode")
	}
}

func TestRotate_MissingSource(t *testing.T) {
	_, err := Rotate("/nonexistent/.env", "/tmp/out.env", RotateOptions{
		Keys: []string{"KEY"},
	})
	if err == nil {
		t.Error("expected error for missing source file")
	}
}

func TestRotate_DefaultsToEmptyValue(t *testing.T) {
	src := writeTempRotateEnv(t, "SECRET=hunter2\n")
	dst := src

	result, err := Rotate(src, dst, RotateOptions{
		Keys: []string{"SECRET"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Records[0].NewValue != "" {
		t.Errorf("expected empty new value, got %q", result.Records[0].NewValue)
	}
}
