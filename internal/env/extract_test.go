package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempExtractEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestExtract_ByKeys(t *testing.T) {
	src := writeTempExtractEnv(t, "FOO=bar\nBAZ=qux\nOTHER=val\n")
	res, err := Extract(src, ExtractOptions{Keys: []string{"FOO", "BAZ"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(res.Records))
	}
	if res.Records[0].Value != "bar" {
		t.Errorf("unexpected value %q", res.Records[0].Value)
	}
}

func TestExtract_MissingKeyMarked(t *testing.T) {
	src := writeTempExtractEnv(t, "FOO=bar\n")
	res, err := Extract(src, ExtractOptions{Keys: []string{"FOO", "NOPE"}})
	if err != nil {
		t.Fatal(err)
	}
	var missing *ExtractRecord
	for i := range res.Records {
		if res.Records[i].Key == "NOPE" {
			missing = &res.Records[i]
		}
	}
	if missing == nil || missing.Found {
		t.Error("expected NOPE to be marked as not found")
	}
}

func TestExtract_ByPrefix(t *testing.T) {
	src := writeTempExtractEnv(t, "APP_HOST=localhost\nAPP_PORT=8080\nDB_URL=postgres\n")
	res, err := Extract(src, ExtractOptions{Prefix: "APP_"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(res.Records))
	}
}

func TestExtract_StripPrefix(t *testing.T) {
	src := writeTempExtractEnv(t, "APP_HOST=localhost\nAPP_PORT=8080\n")
	res, err := Extract(src, ExtractOptions{Prefix: "APP_", Strip: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range res.Records {
		if strings.HasPrefix(r.Key, "APP_") {
			t.Errorf("prefix not stripped from key %q", r.Key)
		}
		if !strings.HasPrefix(r.OriginalKey, "APP_") {
			t.Errorf("original key should retain prefix, got %q", r.OriginalKey)
		}
	}
}

func TestExtract_DryRunNoFile(t *testing.T) {
	src := writeTempExtractEnv(t, "FOO=bar\n")
	out := filepath.Join(t.TempDir(), "out.env")
	_, err := Extract(src, ExtractOptions{Keys: []string{"FOO"}, DryRun: true, Output: out})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Error("dry-run should not create output file")
	}
}

func TestExtract_WritesOutputFile(t *testing.T) {
	src := writeTempExtractEnv(t, "FOO=bar\nBAZ=qux\n")
	out := filepath.Join(t.TempDir(), "out.env")
	_, err := Extract(src, ExtractOptions{Keys: []string{"FOO"}, Output: out})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "FOO=bar") {
		t.Errorf("output file missing expected content, got: %s", data)
	}
	if strings.Contains(string(data), "BAZ") {
		t.Error("output file should not contain non-extracted key BAZ")
	}
}
