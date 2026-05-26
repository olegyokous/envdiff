package diff_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

// sampleResults returns a consistent set of CompareResults for use in report tests.
func sampleResults() []diff.CompareResult {
	return []diff.CompareResult{
		{
			Key:    "DATABASE_URL",
			Status: diff.StatusMissing,
			Values: map[string]string{
				"production": "postgres://prod/db",
				"staging":    "",
			},
		},
		{
			Key:    "API_KEY",
			Status: diff.StatusMismatch,
			Values: map[string]string{
				"production": "prod-secret",
				"staging":    "staging-secret",
			},
		},
		{
			Key:    "PORT",
			Status: diff.StatusMatch,
			Values: map[string]string{
				"production": "8080",
				"staging":    "8080",
			},
		},
	}
}

func TestWriteText_ContainsKeyAndStatus(t *testing.T) {
	var buf bytes.Buffer
	results := sampleResults()

	if err := diff.WriteText(&buf, results); err != nil {
		t.Fatalf("WriteText returned unexpected error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "DATABASE_URL") {
		t.Error("expected output to contain DATABASE_URL")
	}
	if !strings.Contains(output, "MISSING") {
		t.Error("expected output to contain MISSING status")
	}
	if !strings.Contains(output, "MISMATCH") {
		t.Error("expected output to contain MISMATCH status")
	}
	if !strings.Contains(output, "API_KEY") {
		t.Error("expected output to contain API_KEY")
	}
}

func TestWriteText_MatchStatusIncluded(t *testing.T) {
	var buf bytes.Buffer
	results := sampleResults()

	if err := diff.WriteText(&buf, results); err != nil {
		t.Fatalf("WriteText returned unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "PORT") {
		t.Error("expected output to contain PORT key")
	}
	if !strings.Contains(output, "OK") && !strings.Contains(output, "MATCH") {
		t.Error("expected output to contain match status for PORT")
	}
}

func TestWriteText_EmptyResults(t *testing.T) {
	var buf bytes.Buffer

	if err := diff.WriteText(&buf, []diff.CompareResult{}); err != nil {
		t.Fatalf("WriteText returned unexpected error for empty results: %v", err)
	}

	// Output should be empty or just whitespace when there are no results.
	if strings.TrimSpace(buf.String()) != "" {
		t.Errorf("expected empty output for empty results, got: %q", buf.String())
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	results := sampleResults()

	if err := diff.WriteJSON(&buf, results); err != nil {
		t.Fatalf("WriteJSON returned unexpected error: %v", err)
	}

	var parsed []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}

	if len(parsed) != len(results) {
		t.Errorf("expected %d JSON entries, got %d", len(results), len(parsed))
	}
}

func TestWriteJSON_ContainsExpectedFields(t *testing.T) {
	var buf bytes.Buffer
	results := sampleResults()

	if err := diff.WriteJSON(&buf, results); err != nil {
		t.Fatalf("WriteJSON returned unexpected error: %v", err)
	}

	var parsed []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	for _, entry := range parsed {
		if _, ok := entry["key"]; !ok {
			t.Error("expected each JSON entry to have a 'key' field")
		}
		if _, ok := entry["status"]; !ok {
			t.Error("expected each JSON entry to have a 'status' field")
		}
		if _, ok := entry["values"]; !ok {
			t.Error("expected each JSON entry to have a 'values' field")
		}
	}
}
