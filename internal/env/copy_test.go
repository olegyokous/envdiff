package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/envdiff/internal/parser"
)

func writeCopySrc(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "src*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestCopy_CopiesMissingKeys(t *testing.T) {
	src := writeCopySrc(t, "FOO=1\nBAR=2\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	res, err := Copy(src, dst, CopyOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(res.Records))
	}
	for _, r := range res.Records {
		if r.Action != "copied" {
			t.Errorf("key %s: expected copied, got %s", r.Key, r.Action)
		}
	}
	env, _ := parser.ParseFile(dst)
	if env["FOO"] != "1" || env["BAR"] != "2" {
		t.Errorf("destination file content mismatch: %v", env)
	}
}

func TestCopy_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := writeCopySrc(t, "FOO=new\n")
	dst := writeCopySrc(t, "FOO=old\n")

	res, err := Copy(src, dst, CopyOptions{Overwrite: false})
	if err != nil {
		t.Fatal(err)
	}
	if res.Records[0].Action != "skipped" {
		t.Errorf("expected skipped, got %s", res.Records[0].Action)
	}
	env, _ := parser.ParseFile(dst)
	if env["FOO"] != "old" {
		t.Errorf("expected old value preserved, got %s", env["FOO"])
	}
}

func TestCopy_OverwriteReplaces(t *testing.T) {
	src := writeCopySrc(t, "FOO=new\n")
	dst := writeCopySrc(t, "FOO=old\n")

	res, err := Copy(src, dst, CopyOptions{Overwrite: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.Records[0].Action != "overwritten" {
		t.Errorf("expected overwritten, got %s", res.Records[0].Action)
	}
	env, _ := parser.ParseFile(dst)
	if env["FOO"] != "new" {
		t.Errorf("expected new value, got %s", env["FOO"])
	}
}

func TestCopy_DryRunDoesNotWrite(t *testing.T) {
	src := writeCopySrc(t, "FOO=1\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	_, err := Copy(src, dst, CopyOptions{DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		t.Error("expected dst not to exist after dry run")
	}
}

func TestCopy_SelectiveKeys(t *testing.T) {
	src := writeCopySrc(t, "FOO=1\nBAR=2\nBAZ=3\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	res, err := Copy(src, dst, CopyOptions{Keys: []string{"FOO", "BAZ"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(res.Records))
	}
	env, _ := parser.ParseFile(dst)
	if _, ok := env["BAR"]; ok {
		t.Error("BAR should not have been copied")
	}
}
