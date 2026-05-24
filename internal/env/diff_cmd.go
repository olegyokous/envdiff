package env

import (
	"fmt"
	"io"
	"sort"
)

// DiffOptions controls behaviour of the key-level diff between two env maps.
type DiffOptions struct {
	// ShowEqual includes keys that are identical in both envs.
	ShowEqual bool
	// LeftLabel is the display name for the left/source env.
	LeftLabel string
	// RightLabel is the display name for the right/target env.
	RightLabel string
}

// DefaultDiffOptions returns sensible defaults.
func DefaultDiffOptions() DiffOptions {
	return DiffOptions{
		ShowEqual:  false,
		LeftLabel:  "left",
		RightLabel: "right",
	}
}

// KeyDiff represents the comparison result for a single key.
type KeyDiff struct {
	Key        string
	Status     string // "only_left", "only_right", "changed", "equal"
	LeftValue  string
	RightValue string
}

// DiffEnvs compares two env maps and returns per-key differences.
func DiffEnvs(left, right map[string]string, opts DiffOptions) []KeyDiff {
	keys := unionKeys(left, right)
	sort.Strings(keys)

	var results []KeyDiff
	for _, k := range keys {
		lv, inLeft := left[k]
		rv, inRight := right[k]

		var status string
		switch {
		case inLeft && !inRight:
			status = "only_left"
		case !inLeft && inRight:
			status = "only_right"
		case lv == rv:
			status = "equal"
		default:
			status = "changed"
		}

		if status == "equal" && !opts.ShowEqual {
			continue
		}
		results = append(results, KeyDiff{Key: k, Status: status, LeftValue: lv, RightValue: rv})
	}
	return results
}

// WriteDiffText writes a human-readable diff table to w.
func WriteDiffText(w io.Writer, diffs []KeyDiff, opts DiffOptions) {
	if len(diffs) == 0 {
		fmt.Fprintln(w, "No differences found.")
		return
	}
	fmt.Fprintf(w, "%-30s %-12s %-20s %s\n", "KEY", "STATUS", opts.LeftLabel, opts.RightLabel)
	fmt.Fprintf(w, "%s\n", "----------------------------------------------------------------------")
	for _, d := range diffs {
		fmt.Fprintf(w, "%-30s %-12s %-20s %s\n", d.Key, d.Status, d.LeftValue, d.RightValue)
	}
}

// unionKeys returns the set of all keys present in either map.
func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
