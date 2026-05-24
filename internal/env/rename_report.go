package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// RenameRecord captures a single key rename event.
type RenameRecord struct {
	Source string `json:"source"`
	OldKey string `json:"old_key"`
	NewKey string `json:"new_key"`
	Value  string `json:"value"`
	DryRun bool   `json:"dry_run"`
}

// WriteRenameText writes a human-readable rename summary to w.
func WriteRenameText(w io.Writer, r RenameRecord) {
	status := "applied"
	if r.DryRun {
		status = "dry-run"
	}
	fmt.Fprintf(w, "[%s] %s: %s -> %s\n", status, r.Source, r.OldKey, r.NewKey)
}

// WriteRenameJSON writes a JSON-encoded rename record to w.
func WriteRenameJSON(w io.Writer, r RenameRecord) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
