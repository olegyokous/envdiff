package summary_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/summary"
)

func makeResults(statuses ...string) []diff.Result {
	results := make([]diff.Result, 0, len(statuses))
	for i, s := range statuses {
		results = append(results, diff.Result{
			Key:    fmt.Sprintf("KEY_%d", i),
			Status: s,
		})
	}
	return results
}

func TestCompute_AllMatch(t *testing.T) {
	results := []diff.Result{
		{Key: "A", Status: diff.StatusMatch},
		{Key: "B", Status: diff.StatusMatch},
	}
	s := summary.Compute(results)
	if s.Total != 2 || s.Matched != 2 || s.Missing != 0 || s.Mismatch != 0 {
		t.Errorf("unexpected stats: %+v", s)
	}
}

func TestCompute_Mixed(t *testing.T) {
	results := []diff.Result{
		{Key: "A", Status: diff.StatusMatch},
		{Key: "B", Status: diff.StatusMissing},
		{Key: "C", Status: diff.StatusMismatch},
		{Key: "D", Status: diff.StatusMissing},
	}
	s := summary.Compute(results)
	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
	if s.Matched != 1 {
		t.Errorf("expected Matched=1, got %d", s.Matched)
	}
	if s.Missing != 2 {
		t.Errorf("expected Missing=2, got %d", s.Missing)
	}
	if s.Mismatch != 1 {
		t.Errorf("expected Mismatch=1, got %d", s.Mismatch)
	}
}

func TestHasIssues(t *testing.T) {
	clean := summary.Stats{Total: 2, Matched: 2}
	if clean.HasIssues() {
		t.Error("expected no issues for all-match stats")
	}
	dirty := summary.Stats{Total: 3, Matched: 2, Missing: 1}
	if !dirty.HasIssues() {
		t.Error("expected issues when missing > 0")
	}
}

func TestWriteText_OK(t *testing.T) {
	var buf bytes.Buffer
	s := summary.Stats{Total: 3, Matched: 3}
	summary.WriteText(&buf, s)
	out := buf.String()
	if !strings.Contains(out, "Result: OK") {
		t.Errorf("expected OK in output, got: %s", out)
	}
}

func TestWriteText_FAIL(t *testing.T) {
	var buf bytes.Buffer
	s := summary.Stats{Total: 3, Matched: 2, Mismatch: 1}
	summary.WriteText(&buf, s)
	out := buf.String()
	if !strings.Contains(out, "Result: FAIL") {
		t.Errorf("expected FAIL in output, got: %s", out)
	}
	if !strings.Contains(out, "3 keys total") {
		t.Errorf("expected key count in output, got: %s", out)
	}
}
