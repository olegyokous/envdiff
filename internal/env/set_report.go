package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteSetText writes a human-readable summary of a SetResult.
func WriteSetText(w io.Writer, r SetResult) {
	tag := "updated"
	if r.Created {
		tag = "created"
	}
	if r.DryRun {
		tag = "[dry-run] would " + tag
	}

	fmt.Fprintf(w, "file : %s\n", r.File)
	fmt.Fprintf(w, "key  : %s\n", r.Key)
	fmt.Fprintf(w, "value: %s\n", r.Value)
	if !r.Created {
		fmt.Fprintf(w, "prev : %s\n", r.Prev)
	}
	fmt.Fprintf(w, "status: %s\n", tag)
}

// WriteSetJSON writes a JSON representation of a SetResult.
func WriteSetJSON(w io.Writer, r SetResult) error {
	type payload struct {
		File    string `json:"file"`
		Key     string `json:"key"`
		Value   string `json:"value"`
		Prev    string `json:"prev,omitempty"`
		Created bool   `json:"created"`
		DryRun  bool   `json:"dry_run"`
	}
	p := payload{
		File:    r.File,
		Key:     r.Key,
		Value:   r.Value,
		Created: r.Created,
		DryRun:  r.DryRun,
	}
	if !r.Created {
		p.Prev = r.Prev
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}
