package env

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// WriteTrimText writes a human-readable trim report to w.
func WriteTrimText(w io.Writer, result TrimResult) {
	fmt.Fprintf(w, "Trim report for %s\n", result.Source)
	fmt.Fprintf(w, "%s\n", repeatDash(40))

	records := sortedTrimRecords(result.Records)
	for _, r := range records {
		status := "kept"
		if r.Removed {
			status = "removed"
		}
		fmt.Fprintf(w, "  %-30s %s\n", r.Key, status)
	}

	removed := countRemoved(result.Records)
	fmt.Fprintf(w, "%s\n", repeatDash(40))
	fmt.Fprintf(w, "Removed: %d  Kept: %d\n", removed, len(result.Records)-removed)
}

// WriteTrimJSON writes a JSON-encoded trim report to w.
func WriteTrimJSON(w io.Writer, result TrimResult) error {
	type jsonRecord struct {
		Key     string `json:"key"`
		Value   string `json:"value"`
		Removed bool   `json:"removed"`
	}
	type jsonResult struct {
		Source  string       `json:"source"`
		Records []jsonRecord `json:"records"`
	}

	out := jsonResult{Source: result.Source}
	for _, r := range sortedTrimRecords(result.Records) {
		out.Records = append(out.Records, jsonRecord{Key: r.Key, Value: r.Value, Removed: r.Removed})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func sortedTrimRecords(records []TrimRecord) []TrimRecord {
	out := make([]TrimRecord, len(records))
	copy(out, records)
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

func countRemoved(records []TrimRecord) int {
	n := 0
	for _, r := range records {
		if r.Removed {
			n++
		}
	}
	return n
}

func repeatDash(n int) string {
	return fmt.Sprintf("%*s", n, "")[:0] + string(make([]byte, n))
}
