package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// WritePatchText writes a human-readable patch report.
func WritePatchText(w io.Writer, results []PatchResult, source string) {
	fmt.Fprintf(w, "Patch report: %s\n", source)
	fmt.Fprintln(w, strings.Repeat("-", 40))
	if len(results) == 0 {
		fmt.Fprintln(w, "No operations.")
		return
	}
	for _, r := range results {
		status := "OK"
		if !r.Applied {
			status = "SKIP"
		}
		switch r.Op.Action {
		case "set":
			fmt.Fprintf(w, "  [%s] set %s=%s\n", status, r.Op.Key, r.Op.Value)
		case "delete":
			fmt.Fprintf(w, "  [%s] delete %s\n", status, r.Op.Key)
		case "rename":
			fmt.Fprintf(w, "  [%s] rename %s -> %s\n", status, r.Op.Key, r.Op.NewKey)
		default:
			fmt.Fprintf(w, "  [%s] unknown op on %s\n", status, r.Op.Key)
		}
		if r.Reason != "" {
			fmt.Fprintf(w, "         reason: %s\n", r.Reason)
		}
	}
	applied := 0
	for _, r := range results {
		if r.Applied {
			applied++
		}
	}
	fmt.Fprintf(w, "\n%d/%d operations applied.\n", applied, len(results))
}

// WritePatchJSON writes patch results as a JSON array.
func WritePatchJSON(w io.Writer, results []PatchResult) error {
	type jsonRecord struct {
		Action  string `json:"action"`
		Key     string `json:"key"`
		Value   string `json:"value,omitempty"`
		NewKey  string `json:"new_key,omitempty"`
		Applied bool   `json:"applied"`
		Reason  string `json:"reason,omitempty"`
	}
	var records []jsonRecord
	for _, r := range results {
		records = append(records, jsonRecord{
			Action:  r.Op.Action,
			Key:     r.Op.Key,
			Value:   r.Op.Value,
			NewKey:  r.Op.NewKey,
			Applied: r.Applied,
			Reason:  r.Reason,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}
