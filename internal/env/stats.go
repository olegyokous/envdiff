package env

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/parser"
)

// Stats holds aggregate statistics for a single .env file.
type Stats struct {
	File          string         `json:"file"`
	TotalKeys     int            `json:"total_keys"`
	EmptyValues   int            `json:"empty_values"`
	QuotedValues  int            `json:"quoted_values"`
	UniqueValues  int            `json:"unique_values"`
	DuplicateKeys []string       `json:"duplicate_keys,omitempty"`
	PrefixCounts  map[string]int `json:"prefix_counts,omitempty"`
}

// ComputeStats parses the given file and returns a Stats summary.
func ComputeStats(path string) (Stats, error) {
	env, err := parser.ParseFile(path)
	if err != nil {
		return Stats{}, fmt.Errorf("stats: %w", err)
	}

	s := Stats{
		File:         path,
		PrefixCounts: make(map[string]int),
	}

	seen := make(map[string]int)
	valueSeen := make(map[string]bool)

	for k, v := range env {
		seen[k]++
		s.TotalKeys++

		if v == "" {
			s.EmptyValues++
		}

		if len(v) >= 2 && ((v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'')) {
			s.QuotedValues++
		}

		if !valueSeen[v] {
			valueSeen[v] = true
			s.UniqueValues++
		}

		if idx := prefixOf(k); idx != "" {
			s.PrefixCounts[idx]++
		}
	}

	for k, count := range seen {
		if count > 1 {
			s.DuplicateKeys = append(s.DuplicateKeys, k)
		}
	}
	sort.Strings(s.DuplicateKeys)

	return s, nil
}

// prefixOf returns the portion of the key before the first underscore, or "".
func prefixOf(key string) string {
	for i, c := range key {
		if c == '_' && i > 0 {
			return key[:i]
		}
	}
	return ""
}

// WriteStatsText writes a human-readable stats report to w.
func WriteStatsText(w io.Writer, s Stats) {
	fmt.Fprintf(w, "File:           %s\n", s.File)
	fmt.Fprintf(w, "Total keys:     %d\n", s.TotalKeys)
	fmt.Fprintf(w, "Empty values:   %d\n", s.EmptyValues)
	fmt.Fprintf(w, "Quoted values:  %d\n", s.QuotedValues)
	fmt.Fprintf(w, "Unique values:  %d\n", s.UniqueValues)
	if len(s.DuplicateKeys) > 0 {
		fmt.Fprintf(w, "Duplicate keys: %s\n", s.DuplicateKeys)
	}
	if len(s.PrefixCounts) > 0 {
		prefixes := make([]string, 0, len(s.PrefixCounts))
		for p := range s.PrefixCounts {
			prefixes = append(prefixes, p)
		}
		sort.Strings(prefixes)
		fmt.Fprintf(w, "Prefix groups:\n")
		for _, p := range prefixes {
			fmt.Fprintf(w, "  %s: %d\n", p, s.PrefixCounts[p])
		}
	}
}

// WriteStatsJSON writes a JSON-encoded stats report to w.
func WriteStatsJSON(w io.Writer, s Stats) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}
