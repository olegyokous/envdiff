package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteCommentText writes a human-readable report for a CommentResult.
func WriteCommentText(w io.Writer, r CommentResult, src string) {
	switch r.Action {
	case "added":
		fmt.Fprintf(w, "[added]   %s  →  # %s  (%s)\n", r.Key, r.Comment, src)
	case "updated":
		fmt.Fprintf(w, "[updated] %s  →  # %s  (%s)\n", r.Key, r.Comment, src)
	case "removed":
		fmt.Fprintf(w, "[removed] %s  comment removed  (%s)\n", r.Key, src)
	case "not_found":
		fmt.Fprintf(w, "[not_found] %s  key does not exist in %s\n", r.Key, src)
	default:
		fmt.Fprintf(w, "[unknown] %s\n", r.Key)
	}
}

// WriteCommentJSON writes a JSON report for a CommentResult.
func WriteCommentJSON(w io.Writer, r CommentResult, src string) error {
	payload := struct {
		Source  string `json:"source"`
		Key     string `json:"key"`
		Action  string `json:"action"`
		Comment string `json:"comment,omitempty"`
	}{
		Source:  src,
		Key:     r.Key,
		Action:  r.Action,
		Comment: r.Comment,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
