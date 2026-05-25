package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempPatchEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestPatch_SetNewKey(t *testing.T) {
	src := writeTempPatchEnv(t, "FOO=bar\n")
	results, err := Patch(PatchOptions{
		Source: src,
		DryRun: true,
		Ops:    []PatchOp{{Action: "set", Key: "BAZ", Value: "qux"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || !results[0].Applied {
		t.Errorf("expected applied set op, got %+v", results)
	}
}

func TestPatch_DeleteKey(t *testing.T) {
	src := writeTempPatchEnv(t, "FOO=bar\nBAZ=qux\n")
	results, err := Patch(PatchOptions{
		Source: src,
		DryRun: true,
		Ops:    []PatchOp{{Action: "delete", Key: "FOO"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !results[0].Applied {
		t.Errorf("expected delete to be applied")
	}
}

func TestPatch_DeleteMissingKey(t *testing.T) {
	src := writeTempPatchEnv(t, "FOO=bar\n")
	results, err := Patch(PatchOptions{
		Source: src,
		DryRun: true,
		Ops:    []PatchOp{{Action: "delete", Key: "MISSING"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Applied {
		t.Errorf("expected delete to be skipped for missing key")
	}
}

func TestPatch_RenameKey(t *testing.T) {
	src := writeTempPatchEnv(t, "OLD=value\n")
	results, err := Patch(PatchOptions{
		Source: src,
		DryRun: true,
		Ops:    []PatchOp{{Action: "rename", Key: "OLD", NewKey: "NEW"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !results[0].Applied {
		t.Errorf("expected rename to be applied")
	}
}

func TestPatch_RenameConflict(t *testing.T) {
	src := writeTempPatchEnv(t, "OLD=value\nNEW=other\n")
	results, err := Patch(PatchOptions{
		Source: src,
		DryRun: true,
		Ops:    []PatchOp{{Action: "rename", Key: "OLD", NewKey: "NEW"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Applied {
		t.Errorf("expected rename to be skipped due to conflict")
	}
}

func TestPatch_WritesFile(t *testing.T) {
	src := writeTempPatchEnv(t, "FOO=bar\n")
	_, err := Patch(PatchOptions{
		Source: src,
		DryRun: false,
		Ops:    []PatchOp{{Action: "set", Key: "NEW", Value: "val"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(src)
	if !contains(string(data), "NEW=val") {
		t.Errorf("expected NEW=val in output, got: %s", data)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
