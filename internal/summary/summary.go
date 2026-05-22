package summary

import (
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
)

// Stats holds aggregated counts from a diff result set.
type Stats struct {
	Total    int
	Matched  int
	Missing  int
	Mismatch int
}

// Compute calculates summary statistics from a slice of diff results.
func Compute(results []diff.Result) Stats {
	s := Stats{Total: len(results)}
	for _, r := range results {
		switch {
		case r.IsMatch():
			s.Matched++
		case r.IsMissing():
			s.Missing++
		case r.IsMismatch():
			s.Mismatch++
		}
	}
	return s
}

// HasIssues returns true when any missing or mismatched keys exist.
func (s Stats) HasIssues() bool {
	return s.Missing > 0 || s.Mismatch > 0
}

// WriteText writes a human-readable summary line to w.
func WriteText(w io.Writer, s Stats) {
	fmt.Fprintf(w, "Summary: %d keys total — %d matched, %d missing, %d mismatched\n",
		s.Total, s.Matched, s.Missing, s.Mismatch)
	if s.HasIssues() {
		fmt.Fprintln(w, "Result: FAIL")
	} else {
		fmt.Fprintln(w, "Result: OK")
	}
}
