package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// WriteText writes a human-readable diff report to w.
func WriteText(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "✓ No differences found.")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "STATUS\tKEY\tVALUES")
	fmt.Fprintln(tw, "------\t---\t------")

	for _, r := range results {
		valStr := formatValues(r.Values)
		fmt.Fprintf(tw, "%s\t%s\t%s\n", r.Status, r.Key, valStr)
	}
	tw.Flush()
}

// WriteJSON writes a JSON-encoded diff report to w.
func WriteJSON(w io.Writer, results []Result) error {
	if results == nil {
		results = []Result{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

func formatValues(values map[string]string) string {
	if len(values) == 0 {
		return ""
	}
	envs := make([]string, 0, len(values))
	for env := range values {
		envs = append(envs, env)
	}
	sort.Strings(envs)

	out := ""
	for i, env := range envs {
		if i > 0 {
			out += ", "
		}
		out += fmt.Sprintf("%s=%q", env, values[env])
	}
	return out
}
