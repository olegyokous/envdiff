package env

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func sampleSetResult(created bool) SetResult {
	return SetResult{
		File:    "staging.env",
		Key:     "DB_HOST",
		Value:   "localhost",
		Prev:    "old-host",
		Created: created,
		DryRun:  false,
	}
}

func TestWriteSetText_UpdatedKey(t *testing.T) {
	var buf bytes.Buffer
	WriteSetText(&buf, sampleSetResult(false))
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected key in output")
	}
	if !strings.Contains(out, "updated") {
		t.Error("expected 'updated' status")
	}
	if !strings.Contains(out, "old-host") {
		t.Error("expected prev value in output")
	}
}

func TestWriteSetText_CreatedKey(t *testing.T) {
	var buf bytes.Buffer
	WriteSetText(&buf, sampleSetResult(true))
	out := buf.String()
	if !strings.Contains(out, "created") {
		t.Error("expected 'created' status")
	}
	if strings.Contains(out, "prev") {
		t.Error("prev should not appear for new key")
	}
}

func TestWriteSetText_DryRun(t *testing.T) {
	var buf bytes.Buffer
	r := sampleSetResult(false)
	r.DryRun = true
	WriteSetText(&buf, r)
	if !strings.Contains(buf.String(), "dry-run") {
		t.Error("expected dry-run label in output")
	}
}

func TestWriteSetJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteSetJSON(&buf, sampleSetResult(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestWriteSetJSON_ContainsExpectedFields(t *testing.T) {
	var buf bytes.Buffer
	_ = WriteSetJSON(&buf, sampleSetResult(false))
	out := buf.String()
	for _, field := range []string{"file", "key", "value", "prev", "created"} {
		if !strings.Contains(out, field) {
			t.Errorf("expected field %q in JSON output", field)
		}
	}
}
