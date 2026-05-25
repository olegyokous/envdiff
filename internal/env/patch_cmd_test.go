package env

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writePatchCmdEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunPatch_TextOutput(t *testing.T) {
	src := writePatchCmdEnv(t, "FOO=bar\n")
	var buf bytes.Buffer
	err := RunPatch(PatchCmdOptions{
		Source: src,
		DryRun: true,
		Format: "text",
		RawOps: []string{"set:NEW=val"},
	}, &buf)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "set NEW=val") {
		t.Errorf("expected set op in output, got: %s", buf.String())
	}
}

func TestRunPatch_JSONOutput(t *testing.T) {
	src := writePatchCmdEnv(t, "FOO=bar\n")
	var buf bytes.Buffer
	err := RunPatch(PatchCmdOptions{
		Source: src,
		DryRun: true,
		Format: "json",
		RawOps: []string{"delete:FOO"},
	}, &buf)
	if err != nil {
		t.Fatal(err)
	}
	var records []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 1 {
		t.Errorf("expected 1 record, got %d", len(records))
	}
}

func TestRunPatch_InvalidOp(t *testing.T) {
	src := writePatchCmdEnv(t, "FOO=bar\n")
	var buf bytes.Buffer
	err := RunPatch(PatchCmdOptions{
		Source: src,
		RawOps: []string{"badformat"},
	}, &buf)
	if err == nil {
		t.Error("expected error for invalid op format")
	}
}

func TestRunPatch_MissingSource(t *testing.T) {
	var buf bytes.Buffer
	err := RunPatch(PatchCmdOptions{
		Source: "/nonexistent/.env",
		RawOps: []string{"set:K=V"},
	}, &buf)
	if err == nil {
		t.Error("expected error for missing source file")
	}
}

func TestParseRawOps_Rename(t *testing.T) {
	ops, err := parseRawOps([]string{"rename:OLD:NEW"})
	if err != nil {
		t.Fatal(err)
	}
	if ops[0].Key != "OLD" || ops[0].NewKey != "NEW" {
		t.Errorf("unexpected rename op: %+v", ops[0])
	}
}
