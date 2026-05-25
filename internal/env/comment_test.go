package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempCommentEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestComment_AddsInlineComment(t *testing.T) {
	p := writeTempCommentEnv(t, "FOO=bar\nBAZ=qux\n")
	r, err := Comment(p, CommentOptions{Key: "FOO", Comment: "used by auth"})
	if err != nil {
		t.Fatal(err)
	}
	if r.Action != "added" {
		t.Errorf("expected added, got %s", r.Action)
	}
	data, _ := os.ReadFile(p)
	if !strings.Contains(string(data), "FOO=bar # used by auth") {
		t.Errorf("expected inline comment in file, got:\n%s", data)
	}
}

func TestComment_UpdatesExistingComment(t *testing.T) {
	p := writeTempCommentEnv(t, "FOO=bar # old comment\n")
	r, err := Comment(p, CommentOptions{Key: "FOO", Comment: "new comment"})
	if err != nil {
		t.Fatal(err)
	}
	if r.Action != "updated" {
		t.Errorf("expected updated, got %s", r.Action)
	}
	data, _ := os.ReadFile(p)
	if !strings.Contains(string(data), "# new comment") {
		t.Errorf("expected updated comment, got:\n%s", data)
	}
	if strings.Contains(string(data), "old comment") {
		t.Errorf("old comment should be replaced")
	}
}

func TestComment_RemovesComment(t *testing.T) {
	p := writeTempCommentEnv(t, "FOO=bar # remove me\n")
	r, err := Comment(p, CommentOptions{Key: "FOO", Remove: true})
	if err != nil {
		t.Fatal(err)
	}
	if r.Action != "removed" {
		t.Errorf("expected removed, got %s", r.Action)
	}
	data, _ := os.ReadFile(p)
	if strings.Contains(string(data), "#") {
		t.Errorf("expected no comment remaining, got:\n%s", data)
	}
}

func TestComment_KeyNotFound(t *testing.T) {
	p := writeTempCommentEnv(t, "FOO=bar\n")
	r, err := Comment(p, CommentOptions{Key: "MISSING", Comment: "nope"})
	if err != nil {
		t.Fatal(err)
	}
	if r.Action != "not_found" {
		t.Errorf("expected not_found, got %s", r.Action)
	}
}

func TestComment_DryRunDoesNotWrite(t *testing.T) {
	p := writeTempCommentEnv(t, "FOO=bar\n")
	before, _ := os.ReadFile(p)
	_, err := Comment(p, CommentOptions{Key: "FOO", Comment: "dry", DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(p)
	if string(before) != string(after) {
		t.Errorf("dry run should not modify file")
	}
}

func TestComment_MissingFile(t *testing.T) {
	_, err := Comment("/nonexistent/.env", CommentOptions{Key: "X", Comment: "y"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
