package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteRotateText writes a human-readable rotation report to w.
func WriteRotateText(w io.Writer, r RotateResult) {
	fmt.Fprintf(w, "Rotate: %s\n", r.Source)
	if len(r.Records) == 0 {
		fmt.Fprintln(w, "  no keys rotated")
		return
	}
	for _, rec := range r.Records {
		if rec.Rotated {
			fmt.Fprintf(w, "  [rotated] %s  %s -> %s\n", rec.Key, maskValue(rec.OldValue), maskValue(rec.NewValue))
		} else {
			fmt.Fprintf(w, "  [skipped] %s\n", rec.Key)
		}
	}
}

// WriteRotateJSON writes the rotation result as a JSON array to w.
func WriteRotateJSON(w io.Writer, r RotateResult) error {
	type jsonRecord struct {
		Key      string `json:"key"`
		OldValue string `json:"old_value"`
		NewValue string `json:"new_value"`
		Rotated  bool   `json:"rotated"`
	}
	type jsonResult struct {
		Source  string       `json:"source"`
		Records []jsonRecord `json:"records"`
	}

	out := jsonResult{Source: r.Source}
	for _, rec := range r.Records {
		out.Records = append(out.Records, jsonRecord{
			Key:      rec.Key,
			OldValue: rec.OldValue,
			NewValue: rec.NewValue,
			Rotated:  rec.Rotated,
		})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

// maskValue replaces all but the first character with asterisks.
func maskValue(v string) string {
	if len(v) == 0 {
		return "(empty)"
	}
	if len(v) == 1 {
		return "*"
	}
	return string(v[0]) + "***"
}
