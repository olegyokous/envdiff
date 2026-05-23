package template_test

import (
	"testing"

	"github.com/user/envdiff/internal/template"
)

func named(name string, env map[string]string) template.NamedEnv {
	return template.NamedEnv{Name: name, Env: env}
}

func TestMerge_UnionContainsAllKeys(t *testing.T) {
	envs := []template.NamedEnv{
		named("dev", map[string]string{"A": "1", "B": "2"}),
		named("prod", map[string]string{"B": "3", "C": "4"}),
	}
	res := template.Merge(envs)
	for _, k := range []string{"A", "B", "C"} {
		if _, ok := res.Union[k]; !ok {
			t.Errorf("expected key %s in union", k)
		}
	}
}

func TestMerge_FirstValueWins(t *testing.T) {
	envs := []template.NamedEnv{
		named("dev", map[string]string{"X": "first"}),
		named("prod", map[string]string{"X": "second"}),
	}
	res := template.Merge(envs)
	if res.Union["X"] != "first" {
		t.Errorf("expected first value to win, got %s", res.Union["X"])
	}
}

func TestMerge_OnlyInTracksNames(t *testing.T) {
	envs := []template.NamedEnv{
		named("dev", map[string]string{"ONLY_DEV": "1"}),
		named("prod", map[string]string{"ONLY_PROD": "2"}),
	}
	res := template.Merge(envs)
	if len(res.OnlyIn["ONLY_DEV"]) != 1 || res.OnlyIn["ONLY_DEV"][0] != "dev" {
		t.Errorf("unexpected OnlyIn for ONLY_DEV: %v", res.OnlyIn["ONLY_DEV"])
	}
}

func TestMerge_EmptyEnvs(t *testing.T) {
	res := template.Merge(nil)
	if len(res.Union) != 0 {
		t.Error("expected empty union for nil input")
	}
}

func TestUniversalKeys_AllPresent(t *testing.T) {
	envs := []template.NamedEnv{
		named("a", map[string]string{"K1": "1", "K2": "2"}),
		named("b", map[string]string{"K1": "x", "K2": "y"}),
	}
	keys := template.UniversalKeys(envs)
	if len(keys) != 2 {
		t.Errorf("expected 2 universal keys, got %d: %v", len(keys), keys)
	}
}

func TestUniversalKeys_PartialPresence(t *testing.T) {
	envs := []template.NamedEnv{
		named("a", map[string]string{"SHARED": "1", "UNIQUE": "2"}),
		named("b", map[string]string{"SHARED": "3"}),
	}
	keys := template.UniversalKeys(envs)
	if len(keys) != 1 || keys[0] != "SHARED" {
		t.Errorf("expected [SHARED], got %v", keys)
	}
}

func TestUniversalKeys_Empty(t *testing.T) {
	if keys := template.UniversalKeys(nil); len(keys) != 0 {
		t.Errorf("expected empty, got %v", keys)
	}
}
