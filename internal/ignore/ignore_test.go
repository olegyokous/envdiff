package ignore_test

import (
	"os"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/ignore"
)

func writeIgnoreFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envdiff-ignore-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMatch},
		{Key: "SECRET_KEY", Status: diff.StatusMissing},
		{Key: "API_URL", Status: diff.StatusMismatch},
		{Key: "LOG_LEVEL", Status: diff.StatusMatch},
	}
}

func TestNewList_Contains(t *testing.T) {
	l := ignore.NewList([]string{"SECRET_KEY", "LOG_LEVEL"})
	if !l.Contains("SECRET_KEY") {
		t.Error("expected SECRET_KEY to be contained")
	}
	if l.Contains("DB_HOST") {
		t.Error("expected DB_HOST to not be contained")
	}
}

func TestApply_RemovesIgnoredKeys(t *testing.T) {
	l := ignore.NewList([]string{"SECRET_KEY", "LOG_LEVEL"})
	results := l.Apply(sampleResults())
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Key == "SECRET_KEY" || r.Key == "LOG_LEVEL" {
			t.Errorf("key %q should have been filtered out", r.Key)
		}
	}
}

func TestApply_EmptyList(t *testing.T) {
	l := ignore.NewList(nil)
	results := l.Apply(sampleResults())
	if len(results) != len(sampleResults()) {
		t.Errorf("expected all results preserved, got %d", len(results))
	}
}

func TestLoadFile_ParsesKeysAndSkipsComments(t *testing.T) {
	path := writeIgnoreFile(t, "# this is a comment\nSECRET_KEY\n\nAPI_URL\n")
	l, err := ignore.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !l.Contains("SECRET_KEY") {
		t.Error("expected SECRET_KEY")
	}
	if !l.Contains("API_URL") {
		t.Error("expected API_URL")
	}
	if l.Contains("DB_HOST") {
		t.Error("DB_HOST should not be in list")
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	_, err := ignore.LoadFile("/nonexistent/path/.envignore")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
