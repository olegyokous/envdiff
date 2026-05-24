package env

import (
	"fmt"
	"io"
)

// SortCmdOptions bundles CLI-level options for the sort command.
type SortCmdOptions struct {
	Source        string
	Descending    bool
	DryRun        bool
	GroupComments bool
	Quiet         bool
}

// RunSort executes the sort command, writing progress to w.
func RunSort(opts SortCmdOptions, w io.Writer) error {
	if opts.Source == "" {
		return fmt.Errorf("sort: source file is required")
	}

	sopts := SortOptions{
		Descending:    opts.Descending,
		DryRun:        opts.DryRun,
		GroupComments: opts.GroupComments,
	}

	res, err := Sort(opts.Source, sopts)
	if err != nil {
		return err
	}

	if opts.Quiet {
		return nil
	}

	if res.DryRun {
		fmt.Fprintf(w, "dry-run: %s — %d keys (not written)\n", res.Source, res.KeyCount)
	} else {
		fmt.Fprintf(w, "sorted: %s — %d keys\n", res.Source, res.KeyCount)
	}
	return nil
}
