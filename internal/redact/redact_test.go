package redact_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/diff"
	"github.com/your-org/envdiff/internal/redact"
)

func makeResult(key string, values map[string]string) diff.Result {
	return diff.Result{Key: key, Values: values}
}

func TestIsSensitive_MatchesToken(t *testing.T) {
	l := redact.NewList(redact.DefaultSensitivePatterns)
	if !l.IsSensitive("GITHUB_TOKEN") {
		t.Error("expected GITHUB_TOKEN to be sensitive")
	}
}

func TestIsSensitive_MatchesCaseInsensitive(t *testing.T) {
	l := redact.NewList(redact.DefaultSensitivePatterns)
	if !l.IsSensitive("db_password") {
		t.Error("expected db_password to be sensitive")
	}
}

func TestIsSensitive_NotSensitive(t *testing.T) {
	l := redact.NewList(redact.DefaultSensitivePatterns)
	if l.IsSensitive("APP_ENV") {
		t.Error("expected APP_ENV to not be sensitive")
	}
}

func TestApply_MasksSensitiveValues(t *testing.T) {
	l := redact.NewList(redact.DefaultSensitivePatterns)
	results := []diff.Result{
		makeResult("DB_PASSWORD", map[string]string{"prod": "s3cr3t", "dev": "devpass"}),
		makeResult("APP_ENV", map[string]string{"prod": "production", "dev": "development"}),
	}

	out := l.Apply(results)

	for env, val := range out[0].Values {
		if val != redact.Mask {
			t.Errorf("env %s: expected masked value, got %q", env, val)
		}
	}
	if out[1].Values["prod"] != "production" {
		t.Error("non-sensitive value should not be masked")
	}
}

func TestApply_PreservesEmptyValues(t *testing.T) {
	l := redact.NewList(redact.DefaultSensitivePatterns)
	results := []diff.Result{
		makeResult("API_KEY", map[string]string{"prod": "abc123", "dev": ""}),
	}

	out := l.Apply(results)

	if out[0].Values["prod"] != redact.Mask {
		t.Error("non-empty sensitive value should be masked")
	}
	if out[0].Values["dev"] != "" {
		t.Error("empty value should remain empty after masking")
	}
}

func TestApply_CustomPatterns(t *testing.T) {
	l := redact.NewList([]string{"INTERNAL"})
	results := []diff.Result{
		makeResult("INTERNAL_HOST", map[string]string{"prod": "10.0.0.1"}),
		makeResult("PUBLIC_URL", map[string]string{"prod": "https://example.com"}),
	}

	out := l.Apply(results)

	if out[0].Values["prod"] != redact.Mask {
		t.Error("INTERNAL_HOST should be masked")
	}
	if out[1].Values["prod"] == redact.Mask {
		t.Error("PUBLIC_URL should not be masked")
	}
}
