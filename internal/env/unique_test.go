package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempUniqueEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "unique*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestUnique_KeysInOnlyOneFile(t *testing.T) {
	a := writeTempUniqueEnv(t, "SHARED=1\nONLY_A=hello\n")
	b := writeTempUniqueEnv(t, "SHARED=1\nONLY_B=world\n")

	res, err := Unique([]string{a, b}, UniqueOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Records) != 2 {
		t.Fatalf("expected 2 unique records, got %d", len(res.Records))
	}
	keys := map[string]bool{}
	for _, r := range res.Records {
		keys[r.Key] = true
	}
	if !keys["ONLY_A"] || !keys["ONLY_B"] {
		t.Errorf("expected ONLY_A and ONLY_B, got %v", keys)
	}
}

func TestUnique_OnlyInFilter(t *testing.T) {
	a := writeTempUniqueEnv(t, "SHARED=1\nONLY_A=hello\n")
	b := writeTempUniqueEnv(t, "SHARED=1\nONLY_B=world\n")

	res, err := Unique([]string{a, b}, UniqueOptions{OnlyIn: a})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(res.Records))
	}
	if res.Records[0].Key != "ONLY_A" {
		t.Errorf("expected ONLY_A, got %s", res.Records[0].Key)
	}
}

func TestUnique_NoUniqueKeys(t *testing.T) {
	a := writeTempUniqueEnv(t, "KEY=1\n")
	b := writeTempUniqueEnv(t, "KEY=2\n")

	res, err := Unique([]string{a, b}, UniqueOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Records) != 0 {
		t.Errorf("expected 0 records, got %d", len(res.Records))
	}
}

func TestUnique_TooFewSources(t *testing.T) {
	a := writeTempUniqueEnv(t, "KEY=1\n")
	_, err := Unique([]string{a}, UniqueOptions{})
	if err == nil {
		t.Fatal("expected error for single source")
	}
}

func TestUnique_DryRunDoesNotWrite(t *testing.T) {
	a := writeTempUniqueEnv(t, "ONLY_A=1\n")
	b := writeTempUniqueEnv(t, "ONLY_B=2\n")
	out := filepath.Join(t.TempDir(), "out.env")

	_, err := Unique([]string{a, b}, UniqueOptions{DryRun: true, Output: out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Error("expected output file not to be created on dry run")
	}
}

func TestUnique_WritesOutputFile(t *testing.T) {
	a := writeTempUniqueEnv(t, "ONLY_A=val\n")
	b := writeTempUniqueEnv(t, "ONLY_B=other\n")
	out := filepath.Join(t.TempDir(), "out.env")

	_, err := Unique([]string{a, b}, UniqueOptions{Output: out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if len(data) == 0 {
		t.Error("output file is empty")
	}
}
