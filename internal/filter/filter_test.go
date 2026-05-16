package filter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/filter"
)

var sampleResults = []diff.Result{
	{Key: "APP_NAME", Status: "match", Values: map[string]string{"prod": "myapp", "staging": "myapp"}},
	{Key: "APP_SECRET", Status: "mismatch", Values: map[string]string{"prod": "abc", "staging": "xyz"}},
	{Key: "DB_HOST", Status: "missing", Values: map[string]string{"prod": "db.prod", "staging": ""}},
	{Key: "DB_PORT", Status: "match", Values: map[string]string{"prod": "5432", "staging": "5432"}},
	{Key: "REDIS_URL", Status: "missing", Values: map[string]string{"prod": "", "staging": "redis://localhost"}},
}

func TestApply_NoFilter(t *testing.T) {
	got, err := filter.Apply(sampleResults, filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(sampleResults) {
		t.Errorf("expected %d results, got %d", len(sampleResults), len(got))
	}
}

func TestApply_StatusFilter(t *testing.T) {
	got, err := filter.Apply(sampleResults, filter.Options{StatusFilter: "missing"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 missing results, got %d", len(got))
	}
	for _, r := range got {
		if r.Status != "missing" {
			t.Errorf("expected status 'missing', got %q", r.Status)
		}
	}
}

func TestApply_KeyPrefix(t *testing.T) {
	got, err := filter.Apply(sampleResults, filter.Options{KeyPrefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 DB_ results, got %d", len(got))
	}
}

func TestApply_KeyPattern(t *testing.T) {
	got, err := filter.Apply(sampleResults, filter.Options{KeyPattern: "^APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 APP_ results, got %d", len(got))
	}
}

func TestApply_InvalidPattern(t *testing.T) {
	_, err := filter.Apply(sampleResults, filter.Options{KeyPattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}

func TestApply_CombinedFilters(t *testing.T) {
	got, err := filter.Apply(sampleResults, filter.Options{
		StatusFilter: "match",
		KeyPrefix:    "APP_",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 result, got %d", len(got))
	}
	if got[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME, got %q", got[0].Key)
	}
}
