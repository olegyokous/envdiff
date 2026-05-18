package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMatch, Values: map[string]string{"prod": "localhost", "staging": "localhost"}},
		{Key: "API_KEY", Status: diff.StatusMissing, Values: map[string]string{"prod": "secret", "staging": ""}},
	}
}

func TestParseFormat_Valid(t *testing.T) {
	for _, tc := range []struct{ input string; want output.Format }{
		{"text", output.FormatText},
		{"json", output.FormatJSON},
	} {
		got, err := output.ParseFormat(tc.input)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseFormat(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := output.ParseFormat("csv")
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestWrite_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := output.WriterOptions{Format: output.FormatText, Out: &buf}
	if err := output.Write(sampleResults(), opts); err != nil {
		t.Fatalf("Write text: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "DB_HOST") {
		t.Errorf("text output missing DB_HOST: %s", got)
	}
	if !strings.Contains(got, "API_KEY") {
		t.Errorf("text output missing API_KEY: %s", got)
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := output.WriterOptions{Format: output.FormatJSON, Out: &buf}
	if err := output.Write(sampleResults(), opts); err != nil {
		t.Fatalf("Write json: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 results, got %d", len(out))
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := output.WriterOptions{Format: output.Format("xml"), Out: &buf}
	if err := output.Write(sampleResults(), opts); err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := output.DefaultOptions()
	if opts.Format != output.FormatText {
		t.Errorf("DefaultOptions format = %q, want %q", opts.Format, output.FormatText)
	}
	if opts.Out == nil {
		t.Error("DefaultOptions Out should not be nil")
	}
}
