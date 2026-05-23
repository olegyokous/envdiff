package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/envdiff/internal/audit"
	"github.com/your-org/envdiff/internal/diff"
)

func tempLog(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.jsonl")
}

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMatch},
		{Key: "API_KEY", Status: diff.StatusMissing},
	}
}

func TestRecord_CreatesFile(t *testing.T) {
	path := tempLog(t)
	l := audit.NewLogger(path)
	if err := l.Record([]string{"a.env", "b.env"}, sampleResults()); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestRecord_AppendMultiple(t *testing.T) {
	path := tempLog(t)
	l := audit.NewLogger(path)
	for i := 0; i < 3; i++ {
		if err := l.Record([]string{"a.env"}, sampleResults()); err != nil {
			t.Fatalf("Record[%d]: %v", i, err)
		}
	}
	entries, err := audit.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	entries, err := audit.Load("/nonexistent/audit.jsonl")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entries != nil {
		t.Fatalf("expected nil entries, got %v", entries)
	}
}

func TestEntry_HasTimestamp(t *testing.T) {
	path := tempLog(t)
	l := audit.NewLogger(path)
	before := time.Now().UTC().Add(-time.Second)
	if err := l.Record([]string{"x.env"}, sampleResults()); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries, _ := audit.Load(path)
	if entries[0].Timestamp.Before(before) {
		t.Errorf("timestamp too old: %v", entries[0].Timestamp)
	}
}

func TestEntry_SummaryMatchesResults(t *testing.T) {
	path := tempLog(t)
	l := audit.NewLogger(path)
	results := sampleResults()
	_ = l.Record([]string{"a.env"}, results)
	entries, _ := audit.Load(path)
	if entries[0].Summary.Total != len(results) {
		t.Errorf("summary total mismatch: got %d, want %d", entries[0].Summary.Total, len(results))
	}
}
