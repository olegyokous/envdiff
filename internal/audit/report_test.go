package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/your-org/envdiff/internal/audit"
	"github.com/your-org/envdiff/internal/diff"
	"github.com/your-org/envdiff/internal/summary"
)

func makeEntries() []audit.Entry {
	return []audit.Entry{
		{
			Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			Files:     []string{"prod.env", "staging.env"},
			Summary:   summary.Stats{Total: 5, Missing: 1, Mismatch: 2},
			Results:   []diff.Result{{Key: "X", Status: diff.StatusMatch}},
		},
	}
}

func TestWriteText_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	audit.WriteText(&buf, makeEntries())
	if !strings.Contains(buf.String(), "TIMESTAMP") {
		t.Errorf("expected header, got: %s", buf.String())
	}
}

func TestWriteText_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	audit.WriteText(&buf, nil)
	if !strings.Contains(buf.String(), "no audit entries") {
		t.Errorf("expected no-entries message, got: %s", buf.String())
	}
}

func TestWriteJSON_ValidArray(t *testing.T) {
	var buf bytes.Buffer
	if err := audit.WriteJSON(&buf, makeEntries()); err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
}

func TestWriteLast_OnlyLastEntry(t *testing.T) {
	entries := makeEntries()
	entries = append(entries, audit.Entry{
		Timestamp: time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
		Files:     []string{"only.env"},
		Summary:   summary.Stats{Total: 1},
	})
	var buf bytes.Buffer
	audit.WriteLast(&buf, entries)
	if strings.Contains(buf.String(), "2024-06-01") {
		t.Errorf("should only show last entry, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "2024-07-01") {
		t.Errorf("expected last entry date, got: %s", buf.String())
	}
}
