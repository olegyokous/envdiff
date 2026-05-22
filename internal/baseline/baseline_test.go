package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/baseline"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	path := tempPath(t)

	if err := baseline.Save(path, ".env.production", env); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if snap.Source != ".env.production" {
		t.Errorf("Source = %q, want .env.production", snap.Source)
	}
	if snap.Env["FOO"] != "bar" {
		t.Errorf("FOO = %q, want bar", snap.Env["FOO"])
	}
	if snap.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := baseline.Load("/nonexistent/path/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := tempPath(t)
	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := baseline.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestDiff_NoChanges(t *testing.T) {
	snap := &baseline.Snapshot{Env: map[string]string{"A": "1", "B": "2"}}
	current := map[string]string{"A": "1", "B": "2"}
	missing, changed := baseline.Diff(snap, current)
	if len(missing) != 0 || len(changed) != 0 {
		t.Errorf("expected no diff, got missing=%v changed=%v", missing, changed)
	}
}

func TestDiff_MissingKey(t *testing.T) {
	snap := &baseline.Snapshot{Env: map[string]string{"A": "1"}}
	current := map[string]string{"A": "1", "NEW_KEY": "hello"}
	missing, _ := baseline.Diff(snap, current)
	if len(missing) != 1 || missing[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY missing, got %v", missing)
	}
}

func TestDiff_ChangedValue(t *testing.T) {
	snap := &baseline.Snapshot{Env: map[string]string{"PORT": "8080"}}
	current := map[string]string{"PORT": "9090"}
	_, changed := baseline.Diff(snap, current)
	if len(changed) != 1 || changed[0] != "PORT" {
		t.Errorf("expected PORT changed, got %v", changed)
	}
}
