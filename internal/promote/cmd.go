package promote

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// RunConfig holds CLI-level inputs for the promote command.
type RunConfig struct {
	SrcFile  string
	DstFile  string
	Keys     []string
	DryRun   bool
	Overwrite bool
	Out      io.Writer
}

// Execute parses both env files and runs promotion, writing a human-readable
// summary to cfg.Out.
func Execute(cfg RunConfig) error {
	src, err := parser.ParseFile(cfg.SrcFile)
	if err != nil {
		return fmt.Errorf("promote: read source %s: %w", cfg.SrcFile, err)
	}

	dst, err := parser.ParseFile(cfg.DstFile)
	if err != nil {
		return fmt.Errorf("promote: read destination %s: %w", cfg.DstFile, err)
	}

	opts := Options{
		Keys:      cfg.Keys,
		DryRun:    cfg.DryRun,
		Overwrite: cfg.Overwrite,
	}

	results, err := Run(cfg.DstFile, src, dst, opts)
	if err != nil {
		return err
	}

	writeSummary(cfg.Out, results, cfg.DryRun)
	return nil
}

func writeSummary(w io.Writer, results []Result, dryRun bool) {
	if len(results) == 0 {
		fmt.Fprintln(w, "promote: nothing to do")
		return
	}
	for _, r := range results {
		label := strings.ToUpper(r.Status)
		fmt.Fprintf(w, "  [%s] %s\n", label, r.Key)
	}
	promoted := 0
	for _, r := range results {
		if r.Status == "promoted" || r.Status == "dry-run" {
			promoted++
		}
	}
	if dryRun {
		fmt.Fprintf(w, "promote: dry-run — %d key(s) would be promoted\n", promoted)
	} else {
		fmt.Fprintf(w, "promote: %d key(s) promoted\n", promoted)
	}
}
