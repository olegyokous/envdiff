// Package audit records diff runs to an append-only JSONL log file.
package audit

import (
	"encoding/json"
	"os"
	"time"

	"github.com/your-org/envdiff/internal/diff"
	"github.com/your-org/envdiff/internal/summary"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp  time.Time        `json:"timestamp"`
	Files      []string         `json:"files"`
	Summary    summary.Stats    `json:"summary"`
	Results    []diff.Result    `json:"results"`
}

// Logger writes audit entries to a JSONL file.
type Logger struct {
	path string
}

// NewLogger creates a Logger that appends to the given file path.
func NewLogger(path string) *Logger {
	return &Logger{path: path}
}

// Record appends an Entry to the audit log.
func (l *Logger) Record(files []string, results []diff.Result) error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	entry := Entry{
		Timestamp: time.Now().UTC(),
		Files:     files,
		Summary:   summary.Compute(results),
		Results:   results,
	}

	enc := json.NewEncoder(f)
	return enc.Encode(entry)
}

// Load reads all entries from the audit log file.
func Load(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}
