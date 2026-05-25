package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempConvertEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestConvert_ToDotenv(t *testing.T) {
	src := writeTempConvertEnv(t, "FOO=bar\nBAZ=qux\n")
	dst := filepath.Join(t.TempDir(), "out.env")
	r, err := Convert(src, dst, FormatDotenv, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Keys != 2 {
		t.Errorf("expected 2 keys, got %d", r.Keys)
	}
	data, _ := os.ReadFile(dst)
	if !strings.Contains(string(data), "FOO=bar") {
		t.Errorf("output missing FOO=bar: %s", data)
	}
}

func TestConvert_ToExport(t *testing.T) {
	src := writeTempConvertEnv(t, "FOO=bar\n")
	dst := filepath.Join(t.TempDir(), "out.sh")
	r, err := Convert(src, dst, FormatExport, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(r.Output, "export FOO=") {
		t.Errorf("expected export prefix in output: %s", r.Output)
	}
}

func TestConvert_ToJSON(t *testing.T) {
	src := writeTempConvertEnv(t, "KEY=val\n")
	dst := filepath.Join(t.TempDir(), "out.json")
	_, err := Convert(src, dst, FormatJSON, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(dst)
	if !strings.Contains(string(data), "\"KEY\"") {
		t.Errorf("expected JSON key in output: %s", data)
	}
}

func TestConvert_ToYAML(t *testing.T) {
	src := writeTempConvertEnv(t, "KEY=val\n")
	dst := filepath.Join(t.TempDir(), "out.yaml")
	r, err := Convert(src, dst, FormatYAML, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(r.Output, "KEY:") {
		t.Errorf("expected YAML key in output: %s", r.Output)
	}
}

func TestConvert_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempConvertEnv(t, "A=1\n")
	dst := filepath.Join(t.TempDir(), "out.env")
	r, err := Convert(src, dst, FormatDotenv, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.DryRun {
		t.Error("expected DryRun=true")
	}
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		t.Error("expected no file written in dry-run mode")
	}
}

func TestConvert_MissingSource(t *testing.T) {
	_, err := Convert("/no/such/file.env", "/tmp/out.env", FormatDotenv, true)
	if err == nil {
		t.Error("expected error for missing source")
	}
}
