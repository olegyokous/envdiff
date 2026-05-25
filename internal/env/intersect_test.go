package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempIntersectEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestIntersect_KeepsCommonKeys(t *testing.T) {
	src := writeTempIntersectEnv(t, "A=1\nB=2\nC=3\n")
	ref := writeTempIntersectEnv(t, "A=x\nC=y\n")

	opts := DefaultIntersectOptions()
	opts.DryRun = true

	res, err := Intersect(src, ref, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Kept) != 2 {
		t.Errorf("expected 2 kept keys, got %d", len(res.Kept))
	}
	if len(res.Dropped) != 1 || res.Dropped[0] != "B" {
		t.Errorf("expected B dropped, got %v", res.Dropped)
	}
}

func TestIntersect_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempIntersectEnv(t, "A=1\nB=2\n")
	ref := writeTempIntersectEnv(t, "A=x\n")

	original, _ := os.ReadFile(src)

	opts := DefaultIntersectOptions()
	opts.DryRun = true

	_, err := Intersect(src, ref, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	after, _ := os.ReadFile(src)
	if string(original) != string(after) {
		t.Error("dry run should not modify source file")
	}
}

func TestIntersect_WritesToOutputFile(t *testing.T) {
	src := writeTempIntersectEnv(t, "A=1\nB=2\nC=3\n")
	ref := writeTempIntersectEnv(t, "A=x\nC=y\n")
	out := filepath.Join(t.TempDir(), "out.env")

	opts := DefaultIntersectOptions()
	opts.OutputFile = out

	_, err := Intersect(src, ref, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	content := string(data)
	if len(content) == 0 {
		t.Error("output file is empty")
	}
}

func TestIntersect_MissingSource(t *testing.T) {
	ref := writeTempIntersectEnv(t, "A=1\n")
	_, err := Intersect("/nonexistent/.env", ref, DefaultIntersectOptions())
	if err == nil {
		t.Error("expected error for missing source")
	}
}

func TestIntersect_AllDropped(t *testing.T) {
	src := writeTempIntersectEnv(t, "X=1\nY=2\n")
	ref := writeTempIntersectEnv(t, "A=1\nB=2\n")

	opts := DefaultIntersectOptions()
	opts.DryRun = true

	res, err := Intersect(src, ref, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Kept) != 0 {
		t.Errorf("expected no kept keys, got %v", res.Kept)
	}
	if len(res.Dropped) != 2 {
		t.Errorf("expected 2 dropped keys, got %d", len(res.Dropped))
	}
}
