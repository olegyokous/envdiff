package env

import (
	"testing"
)

func TestApply_NoOp(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := Apply(env, DefaultTransformOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("expected unchanged map, got %v", out)
	}
}

func TestApply_KeyPrefix(t *testing.T) {
	env := map[string]string{"NAME": "alice"}
	out, err := Apply(env, TransformOptions{KeyPrefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["APP_NAME"]; !ok {
		t.Errorf("expected key APP_NAME, got %v", out)
	}
}

func TestApply_KeySuffix(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	out, err := Apply(env, TransformOptions{KeySuffix: "_PROD"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["HOST_PROD"]; !ok {
		t.Errorf("expected key HOST_PROD, got %v", out)
	}
}

func TestApply_UpperKeys(t *testing.T) {
	env := map[string]string{"db_host": "localhost"}
	out, err := Apply(env, TransformOptions{UpperKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", out)
	}
}

func TestApply_LowerKeys(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	out, err := Apply(env, TransformOptions{LowerKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["db_host"] != "localhost" {
		t.Errorf("expected db_host=localhost, got %v", out)
	}
}

func TestApply_MutuallyExclusiveKeys(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	_, err := Apply(env, TransformOptions{UpperKeys: true, LowerKeys: true})
	if err == nil {
		t.Error("expected error for conflicting UpperKeys+LowerKeys")
	}
}

func TestApply_ValuePrefixSuffix(t *testing.T) {
	env := map[string]string{"TOKEN": "abc"}
	out, err := Apply(env, TransformOptions{ValuePrefix: "pre_", ValueSuffix: "_suf"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TOKEN"] != "pre_abc_suf" {
		t.Errorf("expected pre_abc_suf, got %q", out["TOKEN"])
	}
}

func TestApply_RenameMap(t *testing.T) {
	env := map[string]string{"OLD_KEY": "value"}
	out, err := Apply(env, TransformOptions{RenameMap: map[string]string{"OLD_KEY": "NEW_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %v", out)
	}
	if _, ok := out["OLD_KEY"]; ok {
		t.Error("OLD_KEY should have been renamed")
	}
}

func TestApply_PrefixThenRename(t *testing.T) {
	env := map[string]string{"HOST": "prod"}
	out, err := Apply(env, TransformOptions{
		KeyPrefix: "APP_",
		RenameMap: map[string]string{"APP_HOST": "SERVICE_HOST"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SERVICE_HOST"] != "prod" {
		t.Errorf("expected SERVICE_HOST=prod, got %v", out)
	}
}
