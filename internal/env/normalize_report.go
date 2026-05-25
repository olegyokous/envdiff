package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteNormalizeText writes a human-readable normalize report to w.
func WriteNormalizeText(w io.Writer, result NormalizeResult, dryRun bool) {
	tag := ""
	if dryRun {
		tag = " [dry-run]"
	}
	fmt.Fprintf(w, "normalize: %s%s\n", result.Source, tag)
	fmt.Fprintln(w, strings.Repeat("-", 40))

	changed := 0
	for _, r := range result.Records {
		if r.Changed {
			changed++
			if r.OldKey != r.Key {
				fmt.Fprintf(w, "  RENAMED  %s -> %s\n", r.OldKey, r.Key)
			}
			if r.OldValue != r.NewValue {
				fmt.Fprintf(w, "  CHANGED  %s: %q -> %q\n", r.Key, r.OldValue, r.NewValue)
			}
		} else {
			fmt.Fprintf(w, "  OK       %s\n", r.Key)
		}
	}
	fmt.Fprintf(w, "\n%d key(s) changed out of %d\n", changed, len(result.Records))
}

// WriteNormalizeJSON writes a JSON array of normalize records to w.
func WriteNormalizeJSON(w io.Writer, result NormalizeResult) error {
	type jsonRecord struct {
		Key      string `json:"key"`
		OldKey   string `json:"old_key,omitempty"`
		OldValue string `json:"old_value"`
		NewValue string `json:"new_value"`
		Changed  bool   `json:"changed"`
	}
	out := make([]jsonRecord, 0, len(result.Records))
	for _, r := range result.Records {
		oldKey := ""
		if r.OldKey != r.Key {
			oldKey = r.OldKey
		}
		out = append(out, jsonRecord{
			Key:      r.Key,
			OldKey:   oldKey,
			OldValue: r.OldValue,
			NewValue: r.NewValue,
			Changed:  r.Changed,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
