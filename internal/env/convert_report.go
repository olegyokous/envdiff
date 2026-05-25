package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteConvertText writes a human-readable conversion summary.
func WriteConvertText(w io.Writer, r *ConvertResult) {
	fmt.Fprintf(w, "convert: %s → %s\n", r.Source, r.Dest)
	fmt.Fprintf(w, "format:  %s\n", r.Format)
	fmt.Fprintf(w, "keys:    %d\n", r.Keys)
	if r.DryRun {
		fmt.Fprintln(w, "mode:    dry-run (no file written)")
		fmt.Fprintln(w, "--- preview ---")
		fmt.Fprint(w, r.Output)
	} else {
		fmt.Fprintf(w, "written: %s\n", r.Dest)
	}
}

// WriteConvertJSON writes the conversion result as JSON.
func WriteConvertJSON(w io.Writer, r *ConvertResult) error {
	payload := map[string]any{
		"source": r.Source,
		"dest":   r.Dest,
		"format": string(r.Format),
		"keys":   r.Keys,
		"dryRun": r.DryRun,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
