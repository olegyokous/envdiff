package env

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func sampleRenameRecord(dryRun bool) RenameRecord {
	return RenameRecord{
		Source: ".env.production",
		OldKey: "DB_HOST",
		NewKey: "DATABASE_HOST",
		Value:  "localhost",
		DryRun: dryRun,
	}
}

func TestWriteRenameText_Applied(t *testing.T) {
	var buf bytes.Buffer
	WriteRenameText(&buf, sampleRenameRecord(false))
	out := buf.String()
	if !strings.Contains(out, "applied") {
		t.Errorf("expected 'applied' in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST -> DATABASE_HOST") {
		t.Errorf("expected key rename in output, got: %s", out)
	}
}

func TestWriteRenameText_DryRun(t *testing.T) {
	var buf bytes.Buffer
	WriteRenameText(&buf, sampleRenameRecord(true))
	out := buf.String()
	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected 'dry-run' in output, got: %s", out)
	}
}

func TestWriteRenameText_ContainsSource(t *testing.T) {
	var buf bytes.Buffer
	WriteRenameText(&buf, sampleRenameRecord(false))
	if !strings.Contains(buf.String(), ".env.production") {
		t.Errorf("expected source filename in output")
	}
}

func TestWriteRenameJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteRenameJSON(&buf, sampleRenameRecord(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got RenameRecord
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.OldKey != "DB_HOST" {
		t.Errorf("expected OldKey=DB_HOST, got %q", got.OldKey)
	}
	if got.NewKey != "DATABASE_HOST" {
		t.Errorf("expected NewKey=DATABASE_HOST, got %q", got.NewKey)
	}
}

func TestWriteRenameJSON_DryRunField(t *testing.T) {
	var buf bytes.Buffer
	_ = WriteRenameJSON(&buf, sampleRenameRecord(true))
	var m map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatal(err)
	}
	if v, ok := m["dry_run"]; !ok || v != true {
		t.Errorf("expected dry_run=true in JSON, got: %v", m)
	}
}
