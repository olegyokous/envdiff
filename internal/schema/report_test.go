package schema_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/schema"
)

func sampleReports() []schema.ViolationReport {
	return []schema.ViolationReport{
		{Env: "production", Violations: []schema.Violation{{Key: "SECRET", Message: "required key is missing or empty"}}, OK: false},
		{Env: "staging", Violations: nil, OK: true},
	}
}

func TestWriteText_ContainsEnvName(t *testing.T) {
	var buf bytes.Buffer
	schema.WriteText(&buf, sampleReports())
	out := buf.String()
	if !strings.Contains(out, "production") {
		t.Error("expected 'production' in text output")
	}
	if !strings.Contains(out, "staging") {
		t.Error("expected 'staging' in text output")
	}
}

func TestWriteText_ShowsViolation(t *testing.T) {
	var buf bytes.Buffer
	schema.WriteText(&buf, sampleReports())
	out := buf.String()
	if !strings.Contains(out, "SECRET") {
		t.Error("expected 'SECRET' in violation output")
	}
}

func TestWriteText_OKEnv(t *testing.T) {
	var buf bytes.Buffer
	schema.WriteText(&buf, sampleReports())
	if !strings.Contains(buf.String(), "schema OK") {
		t.Error("expected 'schema OK' for clean env")
	}
}

func TestWriteJSON_ValidArray(t *testing.T) {
	var buf bytes.Buffer
	if err := schema.WriteJSON(&buf, sampleReports()); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var out []schema.ViolationReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 reports, got %d", len(out))
	}
}

func TestBuild_CreatesReports(t *testing.T) {
	s := &schema.Schema{Fields: []schema.Field{
		{Key: "PORT", Type: schema.TypeInt, Required: true},
	}}
	envs := map[string]map[string]string{
		"prod": {"PORT": "8080"},
		"dev":  {},
	}
	reports := schema.Build(s, envs)
	if len(reports) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(reports))
	}
}

func TestHasViolations_True(t *testing.T) {
	if !schema.HasViolations(sampleReports()) {
		t.Error("expected HasViolations to return true")
	}
}

func TestHasViolations_False(t *testing.T) {
	clean := []schema.ViolationReport{{Env: "prod", OK: true}}
	if schema.HasViolations(clean) {
		t.Error("expected HasViolations to return false")
	}
}
