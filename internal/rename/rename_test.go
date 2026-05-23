package rename

import (
	"testing"
)

func TestDetect_NoRename(t *testing.T) {
	envA := map[string]string{"DB_HOST": "localhost", "PORT": "8080"}
	envB := map[string]string{"DB_HOST": "localhost", "PORT": "8080"}

	got := Detect("dev", envA, "prod", envB)
	if len(got) != 0 {
		t.Errorf("expected no candidates, got %d", len(got))
	}
}

func TestDetect_SimpleRename(t *testing.T) {
	envA := map[string]string{"DB_HOST": "mydb.internal"}
	envB := map[string]string{"DATABASE_HOST": "mydb.internal"}

	got := Detect("dev", envA, "prod", envB)
	if len(got) != 1 {
		t.Fatalf("expected 1 candidate, got %d", len(got))
	}
	if got[0].OldKey != "DB_HOST" {
		t.Errorf("OldKey: got %q, want %q", got[0].OldKey, "DB_HOST")
	}
	if got[0].NewKey != "DATABASE_HOST" {
		t.Errorf("NewKey: got %q, want %q", got[0].NewKey, "DATABASE_HOST")
	}
	if got[0].Value != "mydb.internal" {
		t.Errorf("Value: got %q, want %q", got[0].Value, "mydb.internal")
	}
}

func TestDetect_SkipsCommonValues(t *testing.T) {
	// "true" appears twice in envA so should not be indexed.
	envA := map[string]string{"FEATURE_A": "true", "FEATURE_B": "true"}
	envB := map[string]string{"FEAT_A": "true"}

	got := Detect("dev", envA, "prod", envB)
	if len(got) != 0 {
		t.Errorf("expected no candidates for ambiguous value, got %d", len(got))
	}
}

func TestDetect_SkipsWhenBothKeysExist(t *testing.T) {
	// keyA exists in envB under same name, so not a clean rename.
	envA := map[string]string{"OLD_KEY": "secret123"}
	envB := map[string]string{"OLD_KEY": "other", "NEW_KEY": "secret123"}

	got := Detect("dev", envA, "prod", envB)
	if len(got) != 0 {
		t.Errorf("expected no candidates when old key still exists in envB, got %d", len(got))
	}
}

func TestDetect_MultipleRenames(t *testing.T) {
	envA := map[string]string{
		"SVC_URL":  "https://api.example.com",
		"LOG_PATH": "/var/log/app",
	}
	envB := map[string]string{
		"SERVICE_URL":  "https://api.example.com",
		"LOGGING_PATH": "/var/log/app",
	}

	got := Detect("staging", envA, "prod", envB)
	if len(got) != 2 {
		t.Fatalf("expected 2 candidates, got %d", len(got))
	}
	// Results are sorted by OldKey.
	if got[0].OldKey != "LOG_PATH" {
		t.Errorf("first OldKey: got %q, want %q", got[0].OldKey, "LOG_PATH")
	}
	if got[1].OldKey != "SVC_URL" {
		t.Errorf("second OldKey: got %q, want %q", got[1].OldKey, "SVC_URL")
	}
}

func TestCandidate_String(t *testing.T) {
	c := Candidate{OldKey: "DB_HOST", NewKey: "DATABASE_HOST", Value: "mydb", EnvA: "dev", EnvB: "prod"}
	s := c.String()
	if s == "" {
		t.Error("String() returned empty")
	}
	for _, want := range []string{"DB_HOST", "DATABASE_HOST", "mydb", "dev", "prod"} {
		if !contains(s, want) {
			t.Errorf("String() missing %q in %q", want, s)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
