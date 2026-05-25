package env

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempFlattenEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestFlatten_NoSeparatorKeys(t *testing.T) {
	src := writeTempFlattenEnv(t, "APP_HOST=localhost\nAPP_PORT=5432\n")
	opts := DefaultFlattenOptions()
	opts.Source = src

	r, err := Flatten(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Changed) != 0 {
		t.Errorf("expected no changed keys, got %d", len(r.Changed))
	}
}

func TestFlatten_RenamesSeparatedKeys(t *testing.T) {
	src := writeTempFlattenEnv(t, "APP__DB__HOST=localhost\nAPP__DB__PORT=5432\nPLAIN=value\n")
	opts := DefaultFlattenOptions()
	opts.Source = src

	r, err := Flatten(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Changed) != 2 {
		t.Errorf("expected 2 changed keys, got %d", len(r.Changed))
	}
	if _, ok := r.Flat["APP_DB_HOST"]; !ok {
		t.Error("expected APP_DB_HOST in flat map")
	}
	if _, ok := r.Flat["APP_DB_PORT"]; !ok {
		t.Error("expected APP_DB_PORT in flat map")
	}
	if _, ok := r.Flat["PLAIN"]; !ok {
		t.Error("expected PLAIN unchanged in flat map")
	}
}

func TestFlatten_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempFlattenEnv(t, "APP__KEY=val\n")
	out := filepath.Join(t.TempDir(), "out.env")
	opts := DefaultFlattenOptions()
	opts.Source = src
	opts.Output = out
	opts.DryRun = true

	_, err := Flatten(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Error("expected output file not to be created in dry-run mode")
	}
}

func TestFlatten_WritesOutputFile(t *testing.T) {
	src := writeTempFlattenEnv(t, "APP__KEY=val\n")
	out := filepath.Join(t.TempDir(), "out.env")
	opts := DefaultFlattenOptions()
	opts.Source = src
	opts.Output = out

	_, err := Flatten(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("output file not written: %v", err)
	}
	if !strings.Contains(string(data), "APP_KEY=val") {
		t.Errorf("expected APP_KEY=val in output, got: %s", data)
	}
}

func TestWriteFlattenText_ShowsRenamed(t *testing.T) {
	src := writeTempFlattenEnv(t, "APP__HOST=localhost\n")
	opts := DefaultFlattenOptions()
	opts.Source = src
	r, _ := Flatten(opts)

	var buf bytes.Buffer
	WriteFlattenText(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "APP__HOST") {
		t.Errorf("expected original key in output, got: %s", out)
	}
	if !strings.Contains(out, "APP_HOST") {
		t.Errorf("expected flattened key in output, got: %s", out)
	}
}

func TestWriteFlattenJSON_ValidJSON(t *testing.T) {
	src := writeTempFlattenEnv(t, "APP__HOST=localhost\n")
	opts := DefaultFlattenOptions()
	opts.Source = src
	r, _ := Flatten(opts)

	var buf bytes.Buffer
	if err := WriteFlattenJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["renamed"]; !ok {
		t.Error("expected 'renamed' field in JSON output")
	}
}
