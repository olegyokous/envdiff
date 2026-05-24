package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempDeleteEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestDelete_RemovesKey(t *testing.T) {
	src := writeTempDeleteEnv(t, "FOO=1\nBAR=2\nBAZ=3\n")
	res, err := Delete(src, []string{"BAR"}, DeleteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Records) != 1 || !res.Records[0].Deleted {
		t.Fatalf("expected one deleted record, got %+v", res.Records)
	}
	raw, _ := os.ReadFile(src)
	if strings.Contains(string(raw), "BAR") {
		t.Error("BAR should have been removed from the file")
	}
}

func TestDelete_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempDeleteEnv(t, "FOO=1\nBAR=2\n")
	before, _ := os.ReadFile(src)
	res, err := Delete(src, []string{"FOO"}, DeleteOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.DryRun {
		t.Error("expected DryRun=true in result")
	}
	after, _ := os.ReadFile(src)
	if string(before) != string(after) {
		t.Error("file should not be modified during dry-run")
	}
}

func TestDelete_MissingKeyReturnsError(t *testing.T) {
	src := writeTempDeleteEnv(t, "FOO=1\n")
	_, err := Delete(src, []string{"NOPE"}, DeleteOptions{})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestDelete_MissingKeyForceSkips(t *testing.T) {
	src := writeTempDeleteEnv(t, "FOO=1\n")
	res, err := Delete(src, []string{"NOPE"}, DeleteOptions{Force: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Records) != 1 || !res.Records[0].Missing {
		t.Fatalf("expected missing record, got %+v", res.Records)
	}
}

func TestDelete_PreservesComments(t *testing.T) {
	src := writeTempDeleteEnv(t, "# comment\nFOO=1\nBAR=2\n")
	_, err := Delete(src, []string{"FOO"}, DeleteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	raw, _ := os.ReadFile(src)
	if !strings.Contains(string(raw), "# comment") {
		t.Error("comment should be preserved after delete")
	}
	if strings.Contains(string(raw), "FOO") {
		t.Error("FOO should have been removed")
	}
}

func TestWriteDeleteText_Output(t *testing.T) {
	r := DeleteResult{
		Source: ".env",
		Records: []DeleteRecord{
			{Key: "FOO", Deleted: true},
			{Key: "BAR", Missing: true},
		},
	}
	var sb strings.Builder
	WriteDeleteText(&sb, r)
	out := sb.String()
	if !strings.Contains(out, "DELETED") {
		t.Error("expected DELETED in output")
	}
	if !strings.Contains(out, "MISSING") {
		t.Error("expected MISSING in output")
	}
}
