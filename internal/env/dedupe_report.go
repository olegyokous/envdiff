package env

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// WriteDedupeText writes a human-readable deduplication report.
func WriteDedupeText(w io.Writer, r *DedupeResult) {
	fmt.Fprintf(w, "Source: %s\n", r.Source)
	if r.DryRun {
		fmt.Fprintln(w, "Mode:   dry-run (no changes written)")
	}
	if len(r.Duplicates) == 0 {
		fmt.Fprintln(w, "Status: no duplicate keys found")
		return
	}
	sorted := make([]string, len(r.Duplicates))
	copy(sorted, r.Duplicates)
	sort.Strings(sorted)
	fmt.Fprintf(w, "Duplicates removed (%d):\n", len(sorted))
	for _, k := range sorted {
		fmt.Fprintf(w, "  - %s\n", k)
	}
}

type dedupeJSON struct {
	Source     string   `json:"source"`
	DryRun     bool     `json:"dry_run"`
	Duplicates []string `json:"duplicates"`
	Count      int      `json:"duplicate_count"`
}

// WriteDedupeJSON writes a JSON deduplication report.
func WriteDedupeJSON(w io.Writer, r *DedupeResult) error {
	sorted := make([]string, len(r.Duplicates))
	copy(sorted, r.Duplicates)
	sort.Strings(sorted)
	payload := dedupeJSON{
		Source:     r.Source,
		DryRun:     r.DryRun,
		Duplicates: sorted,
		Count:      len(sorted),
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
