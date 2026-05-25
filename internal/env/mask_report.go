package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteMaskText writes a human-readable mask report.
func WriteMaskText(w io.Writer, result MaskResult) {
	fmt.Fprintf(w, "source: %s\n", result.Source)
	maskedCount := 0
	for _, r := range result.Records {
		if r.Masked {
			maskedCount++
		}
	}
	fmt.Fprintf(w, "keys masked: %d / %d\n\n", maskedCount, len(result.Records))
	for _, r := range result.Records {
		status := "kept"
		if r.Masked {
			status = "masked"
		}
		fmt.Fprintf(w, "  %-30s %s\n", r.Key, status)
	}
}

// WriteMaskJSON writes the mask result as a JSON array.
func WriteMaskJSON(w io.Writer, result MaskResult) error {
	type jsonRecord struct {
		Key    string `json:"key"`
		Masked bool   `json:"masked"`
	}
	type jsonResult struct {
		Source  string       `json:"source"`
		Records []jsonRecord `json:"records"`
	}
	out := jsonResult{Source: result.Source}
	for _, r := range result.Records {
		out.Records = append(out.Records, jsonRecord{Key: r.Key, Masked: r.Masked})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
