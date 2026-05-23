package validate_test

import (
	"regexp"
	"testing"

	"github.com/user/envdiff/internal/validate"
)

func envs(pairs ...map[string]map[string]string) map[string]map[string]string {
	if len(pairs) == 1 {
		return pairs[0]
	}
	return nil
}

func TestCheck_NoViolations(t *testing.T) {
	e := map[string]map[string]string{
		"prod": {"PORT": "8080", "HOST": "example.com"},
	}
	rules := []validate.Rule{
		{Key: "PORT", Pattern: regexp.MustCompile(`^\d+$`)},
		{Key: "HOST", Pattern: nil},
	}
	v := validate.Check(e, rules)
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestCheck_MissingKey(t *testing.T) {
	e := map[string]map[string]string{
		"staging": {"HOST": "staging.example.com"},
	}
	rules := []validate.Rule{
		{Key: "PORT", Pattern: nil},
	}
	v := validate.Check(e, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "PORT" || v[0].Env != "staging" {
		t.Errorf("unexpected violation: %v", v[0])
	}
}

func TestCheck_PatternMismatch(t *testing.T) {
	e := map[string]map[string]string{
		"prod": {"PORT": "not-a-number"},
	}
	rules := []validate.Rule{
		{Key: "PORT", Pattern: regexp.MustCompile(`^\d+$`)},
	}
	v := validate.Check(e, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Message == "" {
		t.Error("expected non-empty message")
	}
}

func TestCheck_MultipleEnvs(t *testing.T) {
	e := map[string]map[string]string{
		"prod":    {"PORT": "8080"},
		"staging": {},
	}
	rules := []validate.Rule{{Key: "PORT", Pattern: nil}}
	v := validate.Check(e, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation (staging missing PORT), got %d", len(v))
	}
}

func TestParseRules_Valid(t *testing.T) {
	raw := map[string]string{"PORT": `^\d+$`, "ENV": `^(prod|staging|dev)$`}
	rules, err := validate.ParseRules(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules))
	}
}

func TestParseRules_InvalidPattern(t *testing.T) {
	raw := map[string]string{"BAD": `[invalid`}
	_, err := validate.ParseRules(raw)
	if err == nil {
		t.Fatal("expected error for invalid regex pattern")
	}
}

func TestViolation_String(t *testing.T) {
	v := validate.Violation{Env: "prod", Key: "PORT", Message: "key is missing"}
	s := v.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}
