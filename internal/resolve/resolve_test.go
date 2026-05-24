package resolve

import (
	"strings"
	"testing"
)

func TestApply_NoRefs(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	res := Apply(env)
	if res.Env["HOST"] != "localhost" {
		t.Errorf("expected localhost, got %q", res.Env["HOST"])
	}
	if len(res.Warnings) != 0 {
		t.Errorf("expected no warnings, got %v", res.Warnings)
	}
}

func TestApply_BraceStyle(t *testing.T) {
	env := map[string]string{
		"HOST":    "localhost",
		"DB_URL":  "postgres://${HOST}:5432/mydb",
	}
	res := Apply(env)
	want := "postgres://localhost:5432/mydb"
	if res.Env["DB_URL"] != want {
		t.Errorf("expected %q, got %q", want, res.Env["DB_URL"])
	}
}

func TestApply_DollarStyle(t *testing.T) {
	env := map[string]string{
		"USER":    "admin",
		"GREETING": "hello $USER",
	}
	res := Apply(env)
	if res.Env["GREETING"] != "hello admin" {
		t.Errorf("unexpected value: %q", res.Env["GREETING"])
	}
}

func TestApply_ChainedRefs(t *testing.T) {
	env := map[string]string{
		"A": "foo",
		"B": "${A}_bar",
		"C": "${B}_baz",
	}
	res := Apply(env)
	if res.Env["C"] != "foo_bar_baz" {
		t.Errorf("chained ref failed: got %q", res.Env["C"])
	}
}

func TestApply_UndefinedRefWarning(t *testing.T) {
	env := map[string]string{
		"URL": "http://${UNDEFINED_HOST}/path",
	}
	res := Apply(env)
	if len(res.Warnings) == 0 {
		t.Fatal("expected a warning for undefined reference")
	}
	if !strings.Contains(res.Warnings[0], "UNDEFINED_HOST") {
		t.Errorf("warning should mention UNDEFINED_HOST, got: %s", res.Warnings[0])
	}
	// Value should remain unresolved.
	if !strings.Contains(res.Env["URL"], "${UNDEFINED_HOST}") {
		t.Errorf("unresolved ref should remain in value, got: %q", res.Env["URL"])
	}
}

func TestApply_OriginalMapUnmodified(t *testing.T) {
	env := map[string]string{
		"BASE": "http://example.com",
		"URL":  "${BASE}/api",
	}
	Apply(env)
	if env["URL"] != "${BASE}/api" {
		t.Error("Apply must not modify the original map")
	}
}
