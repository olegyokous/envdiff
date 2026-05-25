package env

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// WriteFlattenText writes a human-readable summary of the flatten result.
func WriteFlattenText(w io.Writer, r *FlattenResult) {
	fmt.Fprintf(w, "Source : %s\n", r.Source)
	fmt.Fprintf(w, "Total keys : %d\n", len(r.Original))
	fmt.Fprintf(w, "Renamed keys: %d\n", len(r.Changed))

	if len(r.Changed) == 0 {
		fmt.Fprintln(w, "No keys required flattening.")
		return
	}

	sorted := make([]string, len(r.Changed))
	copy(sorted, r.Changed)
	sort.Strings(sorted)

	fmt.Fprintln(w, "\nRenamed:")
	for _, oldKey := range sorted {
		// find new key by checking flat map for a value that matches original
		newKey := findNewKey(r, oldKey)
		fmt.Fprintf(w, "  %s -> %s\n", oldKey, newKey)
	}
}

func findNewKey(r *FlattenResult, oldKey string) string {
	oldVal := r.Original[oldKey]
	for k, v := range r.Flat {
		if v == oldVal && k != oldKey {
			return k
		}
	}
	return oldKey
}

type flattenJSONRecord struct {
	OldKey string `json:"old_key"`
	NewKey string `json:"new_key"`
	Value  string `json:"value"`
}

type flattenJSONReport struct {
	Source  string              `json:"source"`
	Total   int                 `json:"total_keys"`
	Renamed []flattenJSONRecord `json:"renamed"`
}

// WriteFlattenJSON writes a JSON report of the flatten result.
func WriteFlattenJSON(w io.Writer, r *FlattenResult) error {
	records := make([]flattenJSONRecord, 0, len(r.Changed))
	for _, oldKey := range r.Changed {
		newKey := findNewKey(r, oldKey)
		records = append(records, flattenJSONRecord{
			OldKey: oldKey,
			NewKey: newKey,
			Value:  r.Original[oldKey],
		})
	}

	report := flattenJSONReport{
		Source:  r.Source,
		Total:   len(r.Original),
		Renamed: records,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
