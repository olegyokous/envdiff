package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tempMergeOut(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "merged.env")
}

func TestMerge_FirstWins(t *testing.T) {
	a := map[string]string{"KEY": "from_a", "SHARED": "a"}
	b := map[string]string{"KEY": "from_b", "EXTRA": "b"}

	result, err := Merge([]map[string]string{a, b}, DefaultMergeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "from_a" {
		t.Errorf("expected KEY=from_a, got %s", result["KEY"])
	}
	if result["EXTRA"] != "b" {
		t.Errorf("expected EXTRA=b, got %s", result["EXTRA"])
	}
}

func TestMerge_OverwriteReplacesKeys(t *testing.T) {
	a := map[string]string{"KEY": "old"}
	b := map[string]string{"KEY": "new"}

	opts := DefaultMergeOptions()
	opts.Overwrite = true

	result, err := Merge([]map[string]string{a, b}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "new" {
		t.Errorf("expected KEY=new, got %s", result["KEY"])
	}
}

func TestMerge_ExcludesKeys(t *testing.T) {
	a := map[string]string{"KEEP": "yes", "DROP": "no"}

	opts := DefaultMergeOptions()
	opts.ExcludeKeys = map[string]struct{}{"DROP": {}}

	result, err := Merge([]map[string]string{a}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, found := result["DROP"]; found {
		t.Error("expected DROP to be excluded")
	}
	if result["KEEP"] != "yes" {
		t.Errorf("expected KEEP=yes, got %s", result["KEEP"])
	}
}

func TestMerge_WritesFile(t *testing.T) {
	out := tempMergeOut(t)
	a := map[string]string{"APP_NAME": "envdiff", "PORT": "8080"}

	opts := DefaultMergeOptions()
	opts.OutputPath = out
	opts.Header = "merged output"

	_, err := Merge([]map[string]string{a}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "APP_NAME=envdiff") {
		t.Errorf("expected APP_NAME=envdiff in output")
	}
	if !strings.Contains(content, "# merged output") {
		t.Errorf("expected header comment in output")
	}
}

func TestMerge_EmptyEnvs(t *testing.T) {
	result, err := Merge([]map[string]string{}, DefaultMergeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d keys", len(result))
	}
}
