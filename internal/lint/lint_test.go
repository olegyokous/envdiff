package lint_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/lint"
)

func TestCheck_CleanEnv(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "8080",
	}
	warnings := lint.Check("production.env", env)
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %d: %v", len(warnings), warnings)
	}
}

func TestCheck_EmptyValue(t *testing.T) {
	env := map[string]string{
		"SECRET_KEY": "",
	}
	warnings := lint.Check("staging.env", env)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if !strings.Contains(warnings[0].Msg, "empty") {
		t.Errorf("expected empty-value message, got %q", warnings[0].Msg)
	}
	if warnings[0].Key != "SECRET_KEY" {
		t.Errorf("expected key SECRET_KEY, got %q", warnings[0].Key)
	}
}

func TestCheck_WhitespaceOnlyValue(t *testing.T) {
	env := map[string]string{
		"API_TOKEN": "   ",
	}
	warnings := lint.Check("dev.env", env)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if !strings.Contains(warnings[0].Msg, "whitespace") {
		t.Errorf("expected whitespace message, got %q", warnings[0].Msg)
	}
}

func TestCheck_UnusualKeyChar(t *testing.T) {
	env := map[string]string{
		"MY-KEY": "value",
	}
	warnings := lint.Check("test.env", env)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if !strings.Contains(warnings[0].Msg, "unusual character") {
		t.Errorf("expected unusual character message, got %q", warnings[0].Msg)
	}
}

func TestWarning_String(t *testing.T) {
	w := lint.Warning{File: "prod.env", Key: "FOO", Msg: "value is empty"}
	s := w.String()
	if !strings.Contains(s, "prod.env") || !strings.Contains(s, "FOO") || !strings.Contains(s, "value is empty") {
		t.Errorf("String() output missing expected parts: %q", s)
	}
}

func TestCheckAll_MultipleFiles(t *testing.T) {
	envs := map[string]map[string]string{
		"a.env": {"GOOD_KEY": "value", "BAD_KEY": ""},
		"b.env": {"ANOTHER": "ok"},
	}
	warnings := lint.CheckAll(envs)
	if len(warnings) != 1 {
		t.Errorf("expected 1 warning across all files, got %d", len(warnings))
	}
}
