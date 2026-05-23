package template_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/template"
)

func TestGenerate_RedactedOutput(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "secret",
		"APP_ENV": "production",
	}
	opts := template.Options{Redact: true, Header: ""}
	var buf bytes.Buffer
	if err := template.Generate(&buf, env, opts); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, line := range []string{"APP_ENV=", "DB_HOST=", "DB_PASS="} {
		if !strings.Contains(out, line) {
			t.Errorf("expected line %q in output:\n%s", line, out)
		}
	}
	if strings.Contains(out, "secret") || strings.Contains(out, "localhost") {
		t.Error("redacted output should not contain original values")
	}
}

func TestGenerate_PreservesValues(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	opts := template.Options{Redact: false, Header: ""}
	var buf bytes.Buffer
	if err := template.Generate(&buf, env, opts); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "PORT=8080") {
		t.Errorf("expected PORT=8080, got: %s", buf.String())
	}
}

func TestGenerate_SortedKeys(t *testing.T) {
	env := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	opts := template.Options{Redact: false, Header: ""}
	var buf bytes.Buffer
	_ = template.Generate(&buf, env, opts)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if lines[0] != "A_KEY=2" || lines[1] != "M_KEY=3" || lines[2] != "Z_KEY=1" {
		t.Errorf("keys not sorted: %v", lines)
	}
}

func TestGenerate_Header(t *testing.T) {
	var buf bytes.Buffer
	opts := template.DefaultOptions()
	_ = template.Generate(&buf, map[string]string{"X": "y"}, opts)
	if !strings.HasPrefix(buf.String(), "# Auto-generated") {
		t.Errorf("expected header at top, got: %s", buf.String())
	}
}

func TestGenerate_ValueWithSpacesQuoted(t *testing.T) {
	env := map[string]string{"MSG": "hello world"}
	var buf bytes.Buffer
	_ = template.Generate(&buf, env, template.Options{Redact: false})
	if !strings.Contains(buf.String(), `MSG="hello world"`) {
		t.Errorf("expected quoted value, got: %s", buf.String())
	}
}

func TestWriteFile_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.example")
	env := map[string]string{"API_KEY": "abc123"}
	if err := template.WriteFile(path, env, template.DefaultOptions()); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "API_KEY=") {
		t.Errorf("file missing expected key: %s", string(data))
	}
}
