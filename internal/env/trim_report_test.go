package env

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func sampleTrimResult() TrimResult {
	return TrimResult{
		Source: ".env.production",
		Records: []TrimRecord{
			{Key: "API_KEY", Value: "abc123", Removed: false},
			{Key: "UNUSED", Value: "", Removed: true},
			{Key: "DB_HOST", Value: "localhost", Removed: false},
		},
	}
}

func TestWriteTrimText_ContainsSource(t *testing.T) {
	var buf bytes.Buffer
	WriteTrimText(&buf, sampleTrimResult())
	if !strings.Contains(buf.String(), ".env.production") {
		t.Error("expected source file name in output")
	}
}

func TestWriteTrimText_ShowsRemovedAndKept(t *testing.T) {
	var buf bytes.Buffer
	WriteTrimText(&buf, sampleTrimResult())
	out := buf.String()

	if !strings.Contains(out, "removed") {
		t.Error("expected 'removed' status in output")
	}
	if !strings.Contains(out, "kept") {
		t.Error("expected 'kept' status in output")
	}
}

func TestWriteTrimText_CountSummary(t *testing.T) {
	var buf bytes.Buffer
	WriteTrimText(&buf, sampleTrimResult())
	out := buf.String()

	if !strings.Contains(out, "Removed: 1") {
		t.Errorf("expected 'Removed: 1' in summary, got:\n%s", out)
	}
	if !strings.Contains(out, "Kept: 2") {
		t.Errorf("expected 'Kept: 2' in summary, got:\n%s", out)
	}
}

func TestWriteTrimJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteTrimJSON(&buf, sampleTrimResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestWriteTrimJSON_ContainsExpectedFields(t *testing.T) {
	var buf bytes.Buffer
	_ = WriteTrimJSON(&buf, sampleTrimResult())

	if !strings.Contains(buf.String(), "\"source\"") {
		t.Error("expected 'source' field in JSON")
	}
	if !strings.Contains(buf.String(), "\"records\"") {
		t.Error("expected 'records' field in JSON")
	}
	if !strings.Contains(buf.String(), "\"removed\"") {
		t.Error("expected 'removed' field in JSON")
	}
}
