package diff

import (
	"testing"
)

func TestCompare_AllMatch(t *testing.T) {
	envs := map[string]map[string]string{
		"staging": {"DB_HOST": "localhost", "PORT": "5432"},
		"prod":    {"DB_HOST": "localhost", "PORT": "5432"},
	}
	results := Compare(envs)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestCompare_MissingKey(t *testing.T) {
	envs := map[string]map[string]string{
		"staging": {"DB_HOST": "localhost", "SECRET": "abc"},
		"prod":    {"DB_HOST": "localhost"},
	}
	results := Compare(envs)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "SECRET" {
		t.Errorf("expected key SECRET, got %s", results[0].Key)
	}
	if results[0].Status != StatusMissing {
		t.Errorf("expected status missing, got %s", results[0].Status)
	}
}

func TestCompare_MismatchedValue(t *testing.T) {
	envs := map[string]map[string]string{
		"staging": {"LOG_LEVEL": "debug"},
		"prod":    {"LOG_LEVEL": "error"},
	}
	results := Compare(envs)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != StatusMismatch {
		t.Errorf("expected status mismatch, got %s", results[0].Status)
	}
	if results[0].Values["staging"] != "debug" {
		t.Errorf("unexpected staging value: %s", results[0].Values["staging"])
	}
}

func TestCompare_MultipleEnvs(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":     {"APP_ENV": "development", "ONLY_DEV": "1"},
		"staging": {"APP_ENV": "staging"},
		"prod":    {"APP_ENV": "production"},
	}
	results := Compare(envs)
	// ONLY_DEV missing in staging+prod, APP_ENV mismatched
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d: %+v", len(results), results)
	}
}

func TestCompare_EmptyEnvs(t *testing.T) {
	results := Compare(map[string]map[string]string{})
	if results != nil && len(results) != 0 {
		t.Errorf("expected no results for empty input")
	}
}
