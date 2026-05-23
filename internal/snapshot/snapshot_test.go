package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/snapshot"
)

func writeEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestTake_Basic(t *testing.T) {
	dir := t.TempDir()
	f := writeEnv(t, dir, ".env", "FOO=bar\nBAZ=qux\n")
	e, err := snapshot.Take("test", []string{f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Label != "test" {
		t.Errorf("label: got %q, want %q", e.Label, "test")
	}
	state, ok := e.Envs[f]
	if !ok {
		t.Fatal("env state missing")
	}
	if state.Keys["FOO"] != "bar" {
		t.Errorf("FOO: got %q", state.Keys["FOO"])
	}
}

func TestTake_MissingFile(t *testing.T) {
	_, err := snapshot.Take("x", []string{"/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	f := writeEnv(t, dir, ".env", "A=1\n")
	e, _ := snapshot.Take("round", []string{f})
	out := filepath.Join(dir, "snap.json")
	if err := snapshot.Save(out, e); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := snapshot.Load(out)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Label != "round" {
		t.Errorf("label mismatch")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/no/such/file.json")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCompare_DetectsChanges(t *testing.T) {
	dir := t.TempDir()
	f := writeEnv(t, dir, ".env", "FOO=old\nBAR=keep\n")
	old, _ := snapshot.Take("old", []string{f})

	// modify file
	writeEnv(t, dir, ".env", "FOO=new\nBAZ=added\n")
	new, _ := snapshot.Take("new", []string{f})

	diffs := snapshot.Compare(old, new)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	d := diffs[0]
	if len(d.Changed) != 1 || d.Changed[0] != "FOO" {
		t.Errorf("changed: %v", d.Changed)
	}
	if len(d.Added) != 1 || d.Added[0] != "BAZ" {
		t.Errorf("added: %v", d.Added)
	}
	if len(d.Removed) != 1 || d.Removed[0] != "BAR" {
		t.Errorf("removed: %v", d.Removed)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	dir := t.TempDir()
	f := writeEnv(t, dir, ".env", "X=1\n")
	a, _ := snapshot.Take("a", []string{f})
	b, _ := snapshot.Take("b", []string{f})
	diffs := snapshot.Compare(a, b)
	for _, d := range diffs {
		if len(d.Added)+len(d.Removed)+len(d.Changed) > 0 {
			t.Errorf("unexpected diff: %+v", d)
		}
	}
}
