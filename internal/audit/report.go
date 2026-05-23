package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteText writes a human-readable summary of audit entries to w.
func WriteText(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no audit entries found")
		return
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tFILES\tTOTAL\tMISSING\tMISMATCH")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%d\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			len(e.Files),
			e.Summary.Total,
			e.Summary.Missing,
			e.Summary.Mismatch,
		)
	}
	_ = tw.Flush()
}

// WriteJSON writes audit entries as a JSON array to w.
func WriteJSON(w io.Writer, entries []Entry) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

// WriteLast writes only the most recent audit entry in text form.
func WriteLast(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no audit entries found")
		return
	}
	WriteText(w, entries[len(entries)-1:])
}
