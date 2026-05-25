package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeExtractCmdEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunExtract_TextOutput(t *testing.T) {
	src := writeExtractCmdEnv(t, "FOO=bar\nBAZ=qux\n")
	var sb strings.Builder
	err := RunExtract(&sb, ExtractCmdOptions{
		Source: src,
		Keys:   "FOO,BAZ",
		Format: "text",
	})
	if err != nil {
		t.Fatal(err)
	}
	out := sb.String()
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "BAZ") {
		t.Errorf("expected keys in output, got: %s", out)
	}
	if !strings.Contains(out, "found") {
		t.Errorf("expected 'found' status in output, got: %s", out)
	}
}

func TestRunExtract_JSONOutput(t *testing.T) {
	src := writeExtractCmdEnv(t, "FOO=bar\n")
	var sb strings.Builder
	err := RunExtract(&sb, ExtractCmdOptions{
		Source: src,
		Keys:   "FOO",
		Format: "json",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sb.String(), `"key"`) {
		t.Errorf("expected JSON output, got: %s", sb.String())
	}
}

func TestRunExtract_MissingSource(t *testing.T) {
	var sb strings.Builder
	err := RunExtract(&sb, ExtractCmdOptions{
		Source: "/nonexistent/.env",
		Keys:   "FOO",
	})
	if err == nil {
		t.Error("expected error for missing source")
	}
}

func TestRunExtract_NoKeysOrPrefix(t *testing.T) {
	src := writeExtractCmdEnv(t, "FOO=bar\n")
	var sb strings.Builder
	err := RunExtract(&sb, ExtractCmdOptions{Source: src})
	if err == nil {
		t.Error("expected error when neither keys nor prefix provided")
	}
}

func TestRunExtract_PrefixWithStrip(t *testing.T) {
	src := writeExtractCmdEnv(t, "APP_HOST=localhost\nAPP_PORT=8080\nDB_URL=x\n")
	out := filepath.Join(t.TempDir(), "out.env")
	var sb strings.Builder
	err := RunExtract(&sb, ExtractCmdOptions{
		Source: src,
		Prefix: "APP_",
		Strip:  true,
		Output: out,
		Format: "text",
	})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(out)
	if strings.Contains(string(data), "APP_") {
		t.Errorf("prefix should have been stripped in output file, got: %s", data)
	}
}
