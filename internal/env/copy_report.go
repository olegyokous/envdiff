package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteCopyText writes a human-readable summary of a CopyResult.
func WriteCopyText(w io.Writer, r CopyResult) {
	fmt.Fprintf(w, "copy: %s → %s\n", r.Source, r.Destination)
	if len(r.Records) == 0 {
		fmt.Fprintln(w, "  no keys processed")
		return
	}
	for _, rec := range r.Records {
		switch rec.Action {
		case "copied":
			fmt.Fprintf(w, "  + %-30s copied\n", rec.Key)
		case "overwritten":
			fmt.Fprintf(w, "  ~ %-30s overwritten\n", rec.Key)
		case "skipped":
			fmt.Fprintf(w, "  - %-30s skipped (already exists)\n", rec.Key)
		}
	}
}

// WriteCopyJSON writes the CopyResult as a JSON array to w.
func WriteCopyJSON(w io.Writer, r CopyResult) error {
	type jsonOut struct {
		Source      string       `json:"source"`
		Destination string       `json:"destination"`
		Records     []CopyRecord `json:"records"`
	}
	out := jsonOut{
		Source:      r.Source,
		Destination: r.Destination,
		Records:     r.Records,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
