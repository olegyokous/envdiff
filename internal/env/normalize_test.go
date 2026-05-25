package env

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempNormalizeEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestNormalize_UpperKeys(t *testing.T) {
	src := writeTempNormalizeEnv(t, "foo=bar\nbaz=qux\n")
	res, err := Normalize(NormalizeOptions{Source: src, UpperKeys: true, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range res.Records {
		if r.Key != strings.ToUpper(r.OldKey) {
			t.Errorf("expected upper key, got %s", r.Key)
		}
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	src := writeTempNormalizeEnv(t, "KEY=  hello  \n")
	res, err := Normalize(NormalizeOptions{Source: src, TrimValues: true, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Records) == 0 {
		t.Fatal("expected records")
	}
	if res.Records[0].NewValue != "hello" {
		t.Errorf("expected trimmed value, got %q", res.Records[0].NewValue)
	}
}

func TestNormalize_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempNormalizeEnv(t, "key=value\n")
	original, _ := os.ReadFile(src)
	_, err := Normalize(NormalizeOptions{Source: src, UpperKeys: true, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(src)
	if string(original) != string(after) {
		t.Error("dry-run should not modify the file")
	}
}

func TestNormalize_WritesOutput(t *testing.T) {
	src := writeTempNormalizeEnv(t, "key=value\n")
	out := filepath.Join(t.TempDir(), "out.env")
	_, err := Normalize(NormalizeOptions{Source: src, UpperKeys: true, Output: out})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "KEY=") {
		t.Errorf("expected KEY in output, got %s", data)
	}
}

func TestNormalize_ConflictingCaseOptions(t *testing.T) {
	src := writeTempNormalizeEnv(t, "key=val\n")
	_, err := Normalize(NormalizeOptions{Source: src, UpperKeys: true, LowerKeys: true})
	if err == nil {
		t.Error("expected error for conflicting case options")
	}
}

func TestWriteNormalizeJSON_ValidJSON(t *testing.T) {
	result := NormalizeResult{
		Source: ".env",
		Records: []NormalizeRecord{
			{Key: "FOO", OldKey: "foo", OldValue: "bar", NewValue: "bar", Changed: true},
		},
	}
	var buf bytes.Buffer
	if err := WriteNormalizeJSON(&buf, result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"key"`) {
		t.Error("expected key field in JSON output")
	}
}
