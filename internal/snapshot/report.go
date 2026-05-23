package snapshot

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteText writes a human-readable diff report to w.
func WriteText(w io.Writer, diffs []Diff) {
	if len(diffs) == 0 {
		fmt.Fprintln(w, "snapshot: no differences found")
		return
	}
	for _, d := range diffs {
		fmt.Fprintf(w, "file: %s\n", d.File)
		for _, k := range d.Added {
			fmt.Fprintf(w, "  + %s\n", k)
		}
		for _, k := range d.Removed {
			fmt.Fprintf(w, "  - %s\n", k)
		}
		for _, k := range d.Changed {
			fmt.Fprintf(w, "  ~ %s\n", k)
		}
	}
}

// WriteJSON writes diffs as a JSON array to w.
func WriteJSON(w io.Writer, diffs []Diff) error {
	type jsonDiff struct {
		File    string   `json:"file"`
		Added   []string `json:"added"`
		Removed []string `json:"removed"`
		Changed []string `json:"changed"`
	}
	out := make([]jsonDiff, 0, len(diffs))
	for _, d := range diffs {
		out = append(out, jsonDiff{
			File:    d.File,
			Added:   nilSafe(d.Added),
			Removed: nilSafe(d.Removed),
			Changed: nilSafe(d.Changed),
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func nilSafe(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
