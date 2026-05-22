package baseline

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// WriteReport writes a human-readable drift report comparing a snapshot to
// current environment values to w.
func WriteReport(w io.Writer, snap *Snapshot, current map[string]string) error {
	missing, changed := Diff(snap, current)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, "Baseline source:\t%s\n", snap.Source)
	fmt.Fprintf(tw, "Captured at:\t%s\n", snap.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintln(tw, "")

	if len(missing) == 0 && len(changed) == 0 {
		fmt.Fprintln(tw, "Status:\tOK — no drift detected")
		return tw.Flush()
	}

	fmt.Fprintf(tw, "Status:\tDRIFT DETECTED (%d new, %d changed)\n", len(missing), len(changed))
	fmt.Fprintln(tw, "")

	if len(missing) > 0 {
		sort.Strings(missing)
		fmt.Fprintln(tw, "NEW KEYS (not in baseline):")
		for _, k := range missing {
			fmt.Fprintf(tw, "  + %s\t= %s\n", k, current[k])
		}
		fmt.Fprintln(tw, "")
	}

	if len(changed) > 0 {
		sort.Strings(changed)
		fmt.Fprintln(tw, "CHANGED VALUES:")
		for _, k := range changed {
			fmt.Fprintf(tw, "  ~ %s\t%s → %s\n", k, snap.Env[k], current[k])
		}
	}

	return tw.Flush()
}
