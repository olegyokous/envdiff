package diff

import "testing"

func TestResult_IsMatch(t *testing.T) {
	r := Result{Key: "FOO", Status: StatusMatch, Values: map[string]string{"dev": "bar", "prod": "bar"}}
	if !r.IsMatch() {
		t.Error("expected IsMatch() to return true")
	}
	if r.IsMissing() {
		t.Error("expected IsMissing() to return false")
	}
	if r.IsMismatch() {
		t.Error("expected IsMismatch() to return false")
	}
}

func TestResult_IsMissing(t *testing.T) {
	r := Result{Key: "BAR", Status: StatusMissing, Values: map[string]string{"dev": "val", "prod": ""}}
	if !r.IsMissing() {
		t.Error("expected IsMissing() to return true")
	}
	if r.IsMatch() {
		t.Error("expected IsMatch() to return false")
	}
	if r.IsMismatch() {
		t.Error("expected IsMismatch() to return false")
	}
}

func TestResult_IsMismatch(t *testing.T) {
	r := Result{Key: "BAZ", Status: StatusMismatch, Values: map[string]string{"dev": "a", "prod": "b"}}
	if !r.IsMismatch() {
		t.Error("expected IsMismatch() to return true")
	}
	if r.IsMatch() {
		t.Error("expected IsMatch() to return false")
	}
	if r.IsMissing() {
		t.Error("expected IsMissing() to return false")
	}
}

func TestStatusConstants(t *testing.T) {
	if StatusMatch != "match" {
		t.Errorf("unexpected StatusMatch value: %q", StatusMatch)
	}
	if StatusMissing != "missing" {
		t.Errorf("unexpected StatusMissing value: %q", StatusMissing)
	}
	if StatusMismatch != "mismatch" {
		t.Errorf("unexpected StatusMismatch value: %q", StatusMismatch)
	}
}
