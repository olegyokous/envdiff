package env

import (
	"os"
	"strings"
	"testing"
)

func writeTempSortEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestSort_SortsAscending(t *testing.T) {
	src := writeTempSortEnv(t, "ZEBRA=1\nAPPLE=2\nMIDDLE=3\n")
	res, err := Sort(src, DefaultSortOptions())
	if err != nil {
		t.Fatal(err)
	}
	if res.KeyCount != 3 {
		t.Fatalf("expected 3 keys, got %d", res.KeyCount)
	}
	data, _ := os.ReadFile(src)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if lines[0] != "APPLE=2" {
		t.Errorf("first line should be APPLE=2, got %s", lines[0])
	}
	if lines[2] != "ZEBRA=1" {
		t.Errorf("last line should be ZEBRA=1, got %s", lines[2])
	}
}

func TestSort_SortsDescending(t *testing.T) {
	src := writeTempSortEnv(t, "APPLE=2\nZEBRA=1\nMIDDLE=3\n")
	opts := DefaultSortOptions()
	opts.Descending = true
	res, err := Sort(src, opts)
	if err != nil {
		t.Fatal(err)
	}
	if res.KeyCount != 3 {
		t.Fatalf("expected 3 keys, got %d", res.KeyCount)
	}
	data, _ := os.ReadFile(src)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if lines[0] != "ZEBRA=1" {
		t.Errorf("first line should be ZEBRA=1, got %s", lines[0])
	}
}

func TestSort_DryRunDoesNotWrite(t *testing.T) {
	original := "ZEBRA=1\nAPPLE=2\n"
	src := writeTempSortEnv(t, original)
	opts := DefaultSortOptions()
	opts.DryRun = true
	res, err := Sort(src, opts)
	if err != nil {
		t.Fatal(err)
	}
	if !res.DryRun {
		t.Error("expected DryRun=true in result")
	}
	data, _ := os.ReadFile(src)
	if string(data) != original {
		t.Error("dry run should not modify the file")
	}
}

func TestSort_SkipsCommentsAndBlanks(t *testing.T) {
	src := writeTempSortEnv(t, "# header\nZEBRA=1\n\nAPPLE=2\n")
	res, err := Sort(src, DefaultSortOptions())
	if err != nil {
		t.Fatal(err)
	}
	if res.KeyCount != 2 {
		t.Fatalf("expected 2 keys, got %d", res.KeyCount)
	}
}

func TestSort_MissingFile(t *testing.T) {
	_, err := Sort("/nonexistent/.env", DefaultSortOptions())
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSort_GroupComments(t *testing.T) {
	src := writeTempSortEnv(t, "# z comment\nZEBRA=1\n# a comment\nAPPLE=2\n")
	opts := DefaultSortOptions()
	opts.GroupComments = true
	_, err := Sort(src, opts)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(src)
	if !strings.Contains(string(data), "# a comment\nAPPLE=2") {
		t.Errorf("expected comment to precede APPLE=2, got:\n%s", string(data))
	}
}
