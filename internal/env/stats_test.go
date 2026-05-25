package env

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeTempStatsEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempStatsEnv: %v", err)
	}
	return p
}

func TestComputeStats_BasicCounts(t *testing.T) {
	p := writeTempStatsEnv(t, "APP_HOST=localhost\nAPP_PORT=8080\nDB_URL=postgres://localhost\n")
	s, err := ComputeStats(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TotalKeys != 3 {
		t.Errorf("TotalKeys = %d, want 3", s.TotalKeys)
	}
	if s.EmptyValues != 0 {
		t.Errorf("EmptyValues = %d, want 0", s.EmptyValues)
	}
}

func TestComputeStats_EmptyValues(t *testing.T) {
	p := writeTempStatsEnv(t, "FOO=\nBAR=hello\nBAZ=\n")
	s, err := ComputeStats(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.EmptyValues != 2 {
		t.Errorf("EmptyValues = %d, want 2", s.EmptyValues)
	}
}

func TestComputeStats_PrefixGroups(t *testing.T) {
	p := writeTempStatsEnv(t, "APP_HOST=localhost\nAPP_PORT=8080\nDB_URL=postgres://localhost\n")
	s, err := ComputeStats(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.PrefixCounts["APP"] != 2 {
		t.Errorf("PrefixCounts[APP] = %d, want 2", s.PrefixCounts["APP"])
	}
	if s.PrefixCounts["DB"] != 1 {
		t.Errorf("PrefixCounts[DB] = %d, want 1", s.PrefixCounts["DB"])
	}
}

func TestComputeStats_MissingFile(t *testing.T) {
	_, err := ComputeStats("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestWriteStatsText_ContainsFields(t *testing.T) {
	p := writeTempStatsEnv(t, "FOO=bar\nBAZ=\n")
	s, err := ComputeStats(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	WriteStatsText(&buf, s)
	out := buf.String()
	for _, want := range []string{"Total keys", "Empty values", "Unique values"} {
		if !bytes.Contains([]byte(out), []byte(want)) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestWriteStatsJSON_ValidJSON(t *testing.T) {
	p := writeTempStatsEnv(t, "FOO=bar\nBAZ=qux\n")
	s, err := ComputeStats(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	if err := WriteStatsJSON(&buf, s); err != nil {
		t.Fatalf("WriteStatsJSON error: %v", err)
	}
	var decoded Stats
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if decoded.TotalKeys != s.TotalKeys {
		t.Errorf("decoded TotalKeys = %d, want %d", decoded.TotalKeys, s.TotalKeys)
	}
}
