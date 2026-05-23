package promote_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/promote"
)

func writeDst(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env.dst")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRun_PromotesMissingKeys(t *testing.T) {
	dst := writeDst(t, "EXISTING=1\n")
	src := map[string]string{"NEW_KEY": "hello", "EXISTING": "1"}
	dstMap := map[string]string{"EXISTING": "1"}

	results, err := promote.Run(dst, src, dstMap, promote.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "NEW_KEY" || results[0].Status != "promoted" {
		t.Errorf("unexpected results: %+v", results)
	}
}

func TestRun_SkipsExistingWithoutOverwrite(t *testing.T) {
	dst := writeDst(t, "KEY=old\n")
	src := map[string]string{"KEY": "new"}
	dstMap := map[string]string{"KEY": "old"}

	results, err := promote.Run(dst, src, dstMap, promote.Options{Keys: []string{"KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Status != "skipped" {
		t.Errorf("expected skipped, got %+v", results)
	}
}

func TestRun_OverwriteReplaces(t *testing.T) {
	dst := writeDst(t, "KEY=old\n")
	src := map[string]string{"KEY": "new"}
	dstMap := map[string]string{"KEY": "old"}

	results, err := promote.Run(dst, src, dstMap, promote.Options{Keys: []string{"KEY"}, Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Status != "promoted" {
		t.Errorf("expected promoted, got %+v", results)
	}
}

func TestRun_DryRunDoesNotWrite(t *testing.T) {
	dst := writeDst(t, "EXISTING=1\n")
	src := map[string]string{"NEW_KEY": "val"}
	dstMap := map[string]string{"EXISTING": "1"}

	results, err := promote.Run(dst, src, dstMap, promote.Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Status != "dry-run" {
		t.Errorf("expected dry-run, got %+v", results)
	}

	data, _ := os.ReadFile(dst)
	if string(data) != "EXISTING=1\n" {
		t.Errorf("file was modified during dry run")
	}
}

func TestRun_EmptySource(t *testing.T) {
	dst := writeDst(t, "")
	results, err := promote.Run(dst, map[string]string{}, map[string]string{}, promote.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}
