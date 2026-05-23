package drift

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func sampleEntries() []Entry {
	return []Entry{
		{Key: "APP_ENV", Status: StatusMatch, RefValue: "prod", LiveValue: "prod"},
		{Key: "DB_PASS", Status: StatusMissing, RefValue: "secret"},
		{Key: "LOG_LEVEL", Status: StatusDrifted, RefValue: "info", LiveValue: "debug"},
		{Key: "ORPHAN", Status: StatusExtra, LiveValue: "orphan"},
	}
}

func TestWriteText_ContainsKeys(t *testing.T) {
	var buf bytes.Buffer
	WriteText(&buf, sampleEntries())
	out := buf.String()
	for _, key := range []string{"APP_ENV", "DB_PASS", "LOG_LEVEL", "ORPHAN"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %q in text output", key)
		}
	}
}

func TestWriteText_ContainsStatuses(t *testing.T) {
	var buf bytes.Buffer
	WriteText(&buf, sampleEntries())
	out := buf.String()
	for _, token := range []string{"OK", "MISSING", "DRIFTED", "EXTRA"} {
		if !strings.Contains(out, token) {
			t.Errorf("expected status token %q in text output", token)
		}
	}
}

func TestWriteText_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteText(&buf, nil)
	if !strings.Contains(buf.String(), "no entries") {
		t.Error("expected 'no entries' message for empty input")
	}
}

func TestWriteJSON_ValidArray(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteJSON(&buf, sampleEntries()); err != nil {
		t.Fatalf("WriteJSON returned error: %v", err)
	}
	var arr []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &arr); err != nil {
		t.Fatalf("output is not valid JSON array: %v", err)
	}
	if len(arr) != len(sampleEntries()) {
		t.Errorf("expected %d entries, got %d", len(sampleEntries()), len(arr))
	}
}

func TestWriteJSON_ContainsStatusField(t *testing.T) {
	var buf bytes.Buffer
	_ = WriteJSON(&buf, sampleEntries())
	var arr []map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &arr)
	for _, item := range arr {
		if _, ok := item["status"]; !ok {
			t.Errorf("entry missing 'status' field: %v", item)
		}
	}
}
