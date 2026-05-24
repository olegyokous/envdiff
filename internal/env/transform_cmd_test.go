package env

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempTransformSrc(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "src*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestRunTransform_WritesOutput(t *testing.T) {
	src := writeTempTransformSrc(t, "HOST=localhost\nPORT=5432\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	err := RunTransform(TransformCmdOptions{
		Input:  src,
		Output: dst,
		Transform: TransformOptions{KeyPrefix: "APP_"},
	}, os.Stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(dst)
	if !strings.Contains(string(data), "APP_HOST") {
		t.Errorf("expected APP_HOST in output, got:\n%s", data)
	}
}

func TestRunTransform_DryRunNoFile(t *testing.T) {
	src := writeTempTransformSrc(t, "DB=postgres\n")
	dst := filepath.Join(t.TempDir(), "should_not_exist.env")
	var buf bytes.Buffer

	err := RunTransform(TransformCmdOptions{
		Input:  src,
		Output: dst,
		Transform: TransformOptions{UpperKeys: true},
		DryRun: true,
	}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(dst); !os.IsNotExist(statErr) {
		t.Error("dry-run should not write output file")
	}
	if !strings.Contains(buf.String(), "DB") {
		t.Errorf("expected dry-run output to contain key, got: %s", buf.String())
	}
}

func TestRunTransform_MissingInput(t *testing.T) {
	err := RunTransform(TransformCmdOptions{
		Input:  "/no/such/file.env",
		Output: "/tmp/out.env",
	}, os.Stderr)
	if err == nil {
		t.Error("expected error for missing input file")
	}
}

func TestRunTransform_ConflictingCaseOptions(t *testing.T) {
	src := writeTempTransformSrc(t, "KEY=val\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	err := RunTransform(TransformCmdOptions{
		Input:  src,
		Output: dst,
		Transform: TransformOptions{UpperKeys: true, LowerKeys: true},
	}, os.Stderr)
	if err == nil {
		t.Error("expected error for conflicting case options")
	}
}
