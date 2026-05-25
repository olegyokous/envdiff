package env

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func sampleMaskResult() MaskResult {
	return MaskResult{
		Source: ".env.production",
		Records: []MaskRecord{
			{Key: "SECRET", Masked: true},
			{Key: "API_KEY", Masked: true},
			{Key: "PUBLIC", Masked: false},
		},
	}
}

func TestWriteMaskText_ContainsSource(t *testing.T) {
	var buf bytes.Buffer
	WriteMaskText(&buf, sampleMaskResult())
	if !strings.Contains(buf.String(), ".env.production") {
		t.Error("expected source file name in output")
	}
}

func TestWriteMaskText_ShowsMaskedCount(t *testing.T) {
	var buf bytes.Buffer
	WriteMaskText(&buf, sampleMaskResult())
	if !strings.Contains(buf.String(), "keys masked: 2 / 3") {
		t.Errorf("expected masked count summary, got: %s", buf.String())
	}
}

func TestWriteMaskText_ShowsStatus(t *testing.T) {
	var buf bytes.Buffer
	WriteMaskText(&buf, sampleMaskResult())
	out := buf.String()
	if !strings.Contains(out, "masked") {
		t.Error("expected 'masked' status in output")
	}
	if !strings.Contains(out, "kept") {
		t.Error("expected 'kept' status in output")
	}
}

func TestWriteMaskJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteMaskJSON(&buf, sampleMaskResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestWriteMaskJSON_ContainsFields(t *testing.T) {
	var buf bytes.Buffer
	_ = WriteMaskJSON(&buf, sampleMaskResult())
	out := buf.String()
	for _, field := range []string{"source", "records", "key", "masked"} {
		if !strings.Contains(out, field) {
			t.Errorf("expected field %q in JSON output", field)
		}
	}
}
