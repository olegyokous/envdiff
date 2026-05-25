package env

import (
	"encoding/json"
	"strings"
	"testing"
)

func sampleRotateResult() RotateResult {
	return RotateResult{
		Source: ".env.production",
		Records: []RotateRecord{
			{Key: "API_KEY", OldValue: "old", NewValue: "new", Rotated: true},
			{Key: "DB_PASS", OldValue: "pass1", NewValue: "pass2", Rotated: true},
		},
	}
}

func TestWriteRotateText_ContainsHeader(t *testing.T) {
	var buf strings.Builder
	WriteRotateText(&buf, sampleRotateResult())
	if !strings.Contains(buf.String(), ".env.production") {
		t.Error("expected source file name in output")
	}
}

func TestWriteRotateText_ShowsRotatedKeys(t *testing.T) {
	var buf strings.Builder
	WriteRotateText(&buf, sampleRotateResult())
	out := buf.String()
	if !strings.Contains(out, "API_KEY") {
		t.Error("expected API_KEY in output")
	}
	if !strings.Contains(out, "[rotated]") {
		t.Error("expected [rotated] label in output")
	}
}

func TestWriteRotateText_EmptyRecords(t *testing.T) {
	var buf strings.Builder
	WriteRotateText(&buf, RotateResult{Source: "x.env"})
	if !strings.Contains(buf.String(), "no keys rotated") {
		t.Error("expected 'no keys rotated' for empty result")
	}
}

func TestWriteRotateJSON_ValidJSON(t *testing.T) {
	var buf strings.Builder
	if err := WriteRotateJSON(&buf, sampleRotateResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal([]byte(buf.String()), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["source"] != ".env.production" {
		t.Errorf("unexpected source: %v", out["source"])
	}
}

func TestWriteRotateJSON_ContainsRecords(t *testing.T) {
	var buf strings.Builder
	_ = WriteRotateJSON(&buf, sampleRotateResult())
	if !strings.Contains(buf.String(), "API_KEY") {
		t.Error("expected API_KEY in JSON output")
	}
	if !strings.Contains(buf.String(), "rotated") {
		t.Error("expected rotated field in JSON output")
	}
}

func TestMaskValue_Masks(t *testing.T) {
	if maskValue("") != "(empty)" {
		t.Error("expected (empty) for blank value")
	}
	if maskValue("x") != "*" {
		t.Error("expected * for single char")
	}
	masked := maskValue("secret")
	if !strings.HasPrefix(masked, "s") || !strings.Contains(masked, "*") {
		t.Errorf("unexpected mask result: %s", masked)
	}
}
