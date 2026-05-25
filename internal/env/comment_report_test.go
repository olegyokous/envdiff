package env

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func sampleCommentResult(action string) CommentResult {
	return CommentResult{Key: "DB_PASS", Action: action, Comment: "sensitive value"}
}

func TestWriteCommentText_Added(t *testing.T) {
	var buf bytes.Buffer
	WriteCommentText(&buf, sampleCommentResult("added"), ".env")
	out := buf.String()
	if !strings.Contains(out, "[added]") {
		t.Errorf("expected [added] in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected key in output, got: %s", out)
	}
}

func TestWriteCommentText_Updated(t *testing.T) {
	var buf bytes.Buffer
	WriteCommentText(&buf, sampleCommentResult("updated"), ".env")
	if !strings.Contains(buf.String(), "[updated]") {
		t.Errorf("expected [updated] in output")
	}
}

func TestWriteCommentText_Removed(t *testing.T) {
	var buf bytes.Buffer
	WriteCommentText(&buf, sampleCommentResult("removed"), ".env")
	if !strings.Contains(buf.String(), "[removed]") {
		t.Errorf("expected [removed] in output")
	}
}

func TestWriteCommentText_NotFound(t *testing.T) {
	var buf bytes.Buffer
	WriteCommentText(&buf, sampleCommentResult("not_found"), ".env")
	if !strings.Contains(buf.String(), "[not_found]") {
		t.Errorf("expected [not_found] in output")
	}
}

func TestWriteCommentJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteCommentJSON(&buf, sampleCommentResult("added"), ".env.production"); err != nil {
		t.Fatal(err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["key"] != "DB_PASS" {
		t.Errorf("expected key=DB_PASS, got %v", out["key"])
	}
	if out["action"] != "added" {
		t.Errorf("expected action=added, got %v", out["action"])
	}
	if out["source"] != ".env.production" {
		t.Errorf("expected source field, got %v", out["source"])
	}
}

func TestWriteCommentJSON_OmitsEmptyComment(t *testing.T) {
	var buf bytes.Buffer
	r := CommentResult{Key: "X", Action: "removed"}
	if err := WriteCommentJSON(&buf, r, ".env"); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buf.String(), "\"comment\"") {
		t.Errorf("comment field should be omitted when empty")
	}
}
