package exit_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/exit"
)

func TestFromResults_AllMatch(t *testing.T) {
	results := []diff.Result{
		{Key: "FOO", Status: diff.StatusMatch},
		{Key: "BAR", Status: diff.StatusMatch},
	}
	got := exit.FromResults(results)
	if got != exit.CodeOK {
		t.Errorf("expected CodeOK (%d), got %d", exit.CodeOK, got)
	}
}

func TestFromResults_HasMissing(t *testing.T) {
	results := []diff.Result{
		{Key: "FOO", Status: diff.StatusMatch},
		{Key: "BAR", Status: diff.StatusMissing},
	}
	got := exit.FromResults(results)
	if got != exit.CodeMismatch {
		t.Errorf("expected CodeMismatch (%d), got %d", exit.CodeMismatch, got)
	}
}

func TestFromResults_HasMismatch(t *testing.T) {
	results := []diff.Result{
		{Key: "FOO", Status: diff.StatusMismatch},
	}
	got := exit.FromResults(results)
	if got != exit.CodeMismatch {
		t.Errorf("expected CodeMismatch (%d), got %d", exit.CodeMismatch, got)
	}
}

func TestFromResults_Empty(t *testing.T) {
	got := exit.FromResults([]diff.Result{})
	if got != exit.CodeOK {
		t.Errorf("expected CodeOK (%d) for empty results, got %d", exit.CodeOK, got)
	}
}

func TestCodeConstants(t *testing.T) {
	if exit.CodeOK != 0 {
		t.Errorf("CodeOK should be 0, got %d", exit.CodeOK)
	}
	if exit.CodeMismatch != 1 {
		t.Errorf("CodeMismatch should be 1, got %d", exit.CodeMismatch)
	}
	if exit.CodeError != 2 {
		t.Errorf("CodeError should be 2, got %d", exit.CodeError)
	}
}
