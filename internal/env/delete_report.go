package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteDeleteText writes a human-readable summary of a DeleteResult.
func WriteDeleteText(w io.Writer, r DeleteResult) {
	mode := ""
	if r.DryRun {
		mode = " (dry-run)"
	}
	fmt.Fprintf(w, "source: %s%s\n", r.Source, mode)

	if len(r.Records) == 0 {
		fmt.Fprintln(w, "no keys specified")
		return
	}

	for _, rec := range r.Records {
		switch {
		case rec.Missing:
			fmt.Fprintf(w, "  MISSING  %s\n", rec.Key)
		case rec.Deleted:
			fmt.Fprintf(w, "  DELETED  %s\n", rec.Key)
		}
	}
}

type deleteRecordJSON struct {
	Key     string `json:"key"`
	Status  string `json:"status"`
	DryRun  bool   `json:"dry_run"`
}

// WriteDeleteJSON writes a JSON array of delete records.
func WriteDeleteJSON(w io.Writer, r DeleteResult) error {
	var rows []deleteRecordJSON
	for _, rec := range r.Records {
		status := "deleted"
		if rec.Missing {
			status = "missing"
		}
		rows = append(rows, deleteRecordJSON{
			Key:    rec.Key,
			Status: status,
			DryRun: r.DryRun,
		})
	}
	if rows == nil {
		rows = []deleteRecordJSON{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
