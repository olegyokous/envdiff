package drift

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteText writes a human-readable drift report to w.
func WriteText(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no entries to report")
		return
	}
	for _, e := range entries {
		switch e.Status {
		case StatusMatch:
			fmt.Fprintf(w, "  OK       %s\n", e.Key)
		case StatusMissing:
			fmt.Fprintf(w, "  MISSING  %s  (ref: %s)\n", e.Key, e.RefValue)
		case StatusExtra:
			fmt.Fprintf(w, "  EXTRA    %s  (live: %s)\n", e.Key, e.LiveValue)
		case StatusDrifted:
			fmt.Fprintf(w, "  DRIFTED  %s  (ref: %s  live: %s)\n", e.Key, e.RefValue, e.LiveValue)
		}
	}
}

type jsonEntry struct {
	Key       string `json:"key"`
	Status    string `json:"status"`
	RefValue  string `json:"ref_value,omitempty"`
	LiveValue string `json:"live_value,omitempty"`
}

// WriteJSON writes a JSON array of drift entries to w.
func WriteJSON(w io.Writer, entries []Entry) error {
	out := make([]jsonEntry, len(entries))
	for i, e := range entries {
		out[i] = jsonEntry{
			Key:       e.Key,
			Status:    e.Status.String(),
			RefValue:  e.RefValue,
			LiveValue: e.LiveValue,
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
