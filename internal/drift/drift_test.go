package drift

import (
	"testing"
)

func ref() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "db.prod.example.com",
		"LOG_LEVEL": "info",
	}
}

func TestCompare_AllMatch(t *testing.T) {
	live := map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "db.prod.example.com",
		"LOG_LEVEL": "info",
	}
	entries := Compare(ref(), live, false)
	for _, e := range entries {
		if e.Status != StatusMatch {
			t.Errorf("expected match for %s, got %s", e.Key, e.Status)
		}
	}
}

func TestCompare_MissingKey(t *testing.T) {
	live := map[string]string{
		"APP_ENV": "production",
		"DB_HOST": "db.prod.example.com",
	}
	entries := Compare(ref(), live, false)
	var found bool
	for _, e := range entries {
		if e.Key == "LOG_LEVEL" && e.Status == StatusMissing {
			found = true
		}
	}
	if !found {
		t.Error("expected LOG_LEVEL to be reported as missing")
	}
}

func TestCompare_DriftedValue(t *testing.T) {
	live := map[string]string{
		"APP_ENV":  "staging",
		"DB_HOST":  "db.prod.example.com",
		"LOG_LEVEL": "info",
	}
	entries := Compare(ref(), live, false)
	for _, e := range entries {
		if e.Key == "APP_ENV" && e.Status != StatusDrifted {
			t.Errorf("expected APP_ENV drifted, got %s", e.Status)
		}
	}
}

func TestCompare_ExtraKey(t *testing.T) {
	live := map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "db.prod.example.com",
		"LOG_LEVEL": "info",
		"EXTRA_KEY": "surprise",
	}
	entriesNoExtra := Compare(ref(), live, false)
	for _, e := range entriesNoExtra {
		if e.Key == "EXTRA_KEY" {
			t.Error("EXTRA_KEY should not appear when includeExtra=false")
		}
	}
	entriesWithExtra := Compare(ref(), live, true)
	var found bool
	for _, e := range entriesWithExtra {
		if e.Key == "EXTRA_KEY" && e.Status == StatusExtra {
			found = true
		}
	}
	if !found {
		t.Error("expected EXTRA_KEY with StatusExtra when includeExtra=true")
	}
}

func TestHasDrift_True(t *testing.T) {
	entries := []Entry{{Key: "X", Status: StatusMissing}}
	if !HasDrift(entries) {
		t.Error("expected HasDrift to return true")
	}
}

func TestHasDrift_False(t *testing.T) {
	entries := []Entry{{Key: "X", Status: StatusMatch}}
	if HasDrift(entries) {
		t.Error("expected HasDrift to return false")
	}
}

func TestStatusString(t *testing.T) {
	cases := map[Status]string{
		StatusMatch:   "match",
		StatusMissing: "missing",
		StatusExtra:   "extra",
		StatusDrifted: "drifted",
	}
	for s, want := range cases {
		if s.String() != want {
			t.Errorf("Status(%d).String() = %q, want %q", s, s.String(), want)
		}
	}
}
