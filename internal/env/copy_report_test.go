package env

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func sampleCopyResult() CopyResult {
	return CopyResult{
		Source:      "staging.env",
		Destination: "production.env",
		Records: []CopyRecord{
			{Key: "FOO", Value: "bar", Action: "copied"},
			{Key: "BAZ", Value: "qux", Action: "skipped"},
			{Key: "SECRET", Value: "xyz", Action: "overwritten"},
		},
	}
}

func TestWriteCopyText_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	WriteCopyText(&buf, sampleCopyResult())
	if !strings.Contains(buf.String(), "staging.env") {
		t.Error("expected source name in output")
	}
	if !strings.Contains(buf.String(), "production.env") {
		t.Error("expected destination name in output")
	}
}

func TestWriteCopyText_ShowsActions(t *testing.T) {
	var buf bytes.Buffer
	WriteCopyText(&buf, sampleCopyResult())
	out := buf.String()
	if !strings.Contains(out, "copied") {
		t.Error("expected 'copied' in output")
	}
	if !strings.Contains(out, "skipped") {
		t.Error("expected 'skipped' in output")
	}
	if !strings.Contains(out, "overwritten") {
		t.Error("expected 'overwritten' in output")
	}
}

func TestWriteCopyText_EmptyRecords(t *testing.T) {
	var buf bytes.Buffer
	WriteCopyText(&buf, CopyResult{Source: "a.env", Destination: "b.env"})
	if !strings.Contains(buf.String(), "no keys processed") {
		t.Error("expected 'no keys processed' message")
	}
}

func TestWriteCopyJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteCopyJSON(&buf, sampleCopyResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestWriteCopyJSON_ContainsFields(t *testing.T) {
	var buf bytes.Buffer
	WriteCopyJSON(&buf, sampleCopyResult())
	out := buf.String()
	for _, field := range []string{"source", "destination", "records", "action", "key"} {
		if !strings.Contains(out, field) {
			t.Errorf("expected field %q in JSON output", field)
		}
	}
}
