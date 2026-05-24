package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// grepMatchJSON is the JSON-serialisable form of a GrepMatch.
type grepMatchJSON struct {
	File  string `json:"file"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// WriteGrepJSON writes matches as a JSON array.
func WriteGrepJSON(w io.Writer, matches []GrepMatch) error {
	out := make([]grepMatchJSON, len(matches))
	for i, m := range matches {
		out[i] = grepMatchJSON{File: m.File, Key: m.Key, Value: m.Value}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return fmt.Errorf("grep: json encode: %w", err)
	}
	return nil
}

// WriteGrepSummary prints a one-line summary of the grep results.
func WriteGrepSummary(w io.Writer, matches []GrepMatch) {
	if len(matches) == 0 {
		fmt.Fprintln(w, "grep: 0 matches")
		return
	}
	files := map[string]struct{}{}
	for _, m := range matches {
		files[m.File] = struct{}{}
	}
	fmt.Fprintf(w, "grep: %d match(es) across %d file(s)\n", len(matches), len(files))
}
