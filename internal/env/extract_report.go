package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteExtractText writes a human-readable extract report.
func WriteExtractText(w io.Writer, res ExtractResult) {
	fmt.Fprintf(w, "Source: %s\n", res.Source)
	fmt.Fprintf(w, "%-32s %s\n", "KEY", "STATUS")
	fmt.Fprintln(w, repeatDash(48))
	for _, r := range res.Records {
		status := "found"
		if !r.Found {
			status = "missing"
		}
		key := r.Key
		if r.OriginalKey != r.Key {
			key = fmt.Sprintf("%s (was %s)", r.Key, r.OriginalKey)
		}
		fmt.Fprintf(w, "%-32s %s\n", key, status)
	}
	found := 0
	for _, r := range res.Records {
		if r.Found {
			found++
		}
	}
	fmt.Fprintf(w, "\n%d/%d keys extracted\n", found, len(res.Records))
}

// WriteExtractJSON writes a JSON array of extract records.
func WriteExtractJSON(w io.Writer, res ExtractResult) error {
	type jsonRecord struct {
		Key         string `json:"key"`
		OriginalKey string `json:"original_key,omitempty"`
		Value       string `json:"value,omitempty"`
		Found       bool   `json:"found"`
		Source      string `json:"source"`
	}
	out := make([]jsonRecord, 0, len(res.Records))
	for _, r := range res.Records {
		orig := ""
		if r.OriginalKey != r.Key {
			orig = r.OriginalKey
		}
		out = append(out, jsonRecord{
			Key: r.Key, OriginalKey: orig,
			Value: r.Value, Found: r.Found,
			Source: res.Source,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
