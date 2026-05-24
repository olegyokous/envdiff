package env

import (
	"bytes"
	"strings"
	"testing"
)

func TestDiffEnvs_NoChanges(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1", "B": "2"}
	opts := DefaultDiffOptions()

	diffs := DiffEnvs(left, right, opts)
	if len(diffs) != 0 {
		t.Fatalf("expected 0 diffs, got %d", len(diffs))
	}
}

func TestDiffEnvs_ShowEqual(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"A": "1"}
	opts := DefaultDiffOptions()
	opts.ShowEqual = true

	diffs := DiffEnvs(left, right, opts)
	if len(diffs) != 1 || diffs[0].Status != "equal" {
		t.Fatalf("expected 1 equal diff, got %+v", diffs)
	}
}

func TestDiffEnvs_OnlyLeft(t *testing.T) {
	left := map[string]string{"A": "1", "X": "extra"}
	right := map[string]string{"A": "1"}
	opts := DefaultDiffOptions()

	diffs := DiffEnvs(left, right, opts)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Status != "only_left" || diffs[0].Key != "X" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
}

func TestDiffEnvs_OnlyRight(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"A": "1", "NEW": "val"}
	opts := DefaultDiffOptions()

	diffs := DiffEnvs(left, right, opts)
	if len(diffs) != 1 || diffs[0].Status != "only_right" {
		t.Fatalf("expected only_right diff, got %+v", diffs)
	}
}

func TestDiffEnvs_Changed(t *testing.T) {
	left := map[string]string{"DB": "old"}
	right := map[string]string{"DB": "new"}
	opts := DefaultDiffOptions()

	diffs := DiffEnvs(left, right, opts)
	if len(diffs) != 1 || diffs[0].Status != "changed" {
		t.Fatalf("expected changed diff, got %+v", diffs)
	}
	if diffs[0].LeftValue != "old" || diffs[0].RightValue != "new" {
		t.Errorf("wrong values: %+v", diffs[0])
	}
}

func TestWriteDiffText_NoDiffs(t *testing.T) {
	var buf bytes.Buffer
	WriteDiffText(&buf, nil, DefaultDiffOptions())
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestWriteDiffText_ContainsKey(t *testing.T) {
	diffs := []KeyDiff{
		{Key: "SECRET", Status: "changed", LeftValue: "abc", RightValue: "xyz"},
	}
	var buf bytes.Buffer
	WriteDiffText(&buf, diffs, DefaultDiffOptions())
	out := buf.String()
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected key SECRET in output: %s", out)
	}
	if !strings.Contains(out, "changed") {
		t.Errorf("expected status 'changed' in output: %s", out)
	}
}
