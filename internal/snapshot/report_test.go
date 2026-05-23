package snapshot_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/snapshot"
)

func sampleDiffs() []snapshot.Diff {
	return []snapshot.Diff{
		{
			File:    ".env.production",
			Added:   []string{"NEW_KEY"},
			Removed: []string{"OLD_KEY"},
			Changed: []string{"CHANGED_KEY"},
		},
	}
}

func TestWriteText_ContainsFile(t *testing.T) {
	var buf bytes.Buffer
	snapshot.WriteText(&buf, sampleDiffs())
	if !strings.Contains(buf.String(), ".env.production") {
		t.Errorf("expected file name in output: %s", buf.String())
	}
}

func TestWriteText_ShowsAddedRemovedChanged(t *testing.T) {
	var buf bytes.Buffer
	snapshot.WriteText(&buf, sampleDiffs())
	out := buf.String()
	for _, want := range []string{"+ NEW_KEY", "- OLD_KEY", "~ CHANGED_KEY"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in output", want)
		}
	}
}

func TestWriteText_NoDiffs(t *testing.T) {
	var buf bytes.Buffer
	snapshot.WriteText(&buf, nil)
	if !strings.Contains(buf.String(), "no differences") {
		t.Errorf("expected no-diff message")
	}
}

func TestWriteJSON_ValidArray(t *testing.T) {
	var buf bytes.Buffer
	if err := snapshot.WriteJSON(&buf, sampleDiffs()); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var arr []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &arr); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(arr) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(arr))
	}
}

func TestWriteJSON_ContainsExpectedFields(t *testing.T) {
	var buf bytes.Buffer
	snapshot.WriteJSON(&buf, sampleDiffs())
	out := buf.String()
	for _, field := range []string{"file", "added", "removed", "changed"} {
		if !strings.Contains(out, field) {
			t.Errorf("missing field %q in JSON", field)
		}
	}
}

func TestWriteJSON_EmptySlicesNotNull(t *testing.T) {
	var buf bytes.Buffer
	snapshot.WriteJSON(&buf, []snapshot.Diff{{File: "x", Added: nil}})
	if strings.Contains(buf.String(), "null") {
		t.Error("null found in JSON; expected empty arrays")
	}
}
