package env

import (
	"fmt"
	"io"
	"strings"
)

// ExtractCmdOptions are the CLI-level options for RunExtract.
type ExtractCmdOptions struct {
	Source  string
	Keys    string // comma-separated
	Prefix  string
	Strip   bool
	DryRun  bool
	Output  string
	Format  string // "text" | "json"
}

// RunExtract is the entry-point called by the CLI subcommand.
func RunExtract(w io.Writer, opts ExtractCmdOptions) error {
	if opts.Source == "" {
		return fmt.Errorf("extract: source file is required")
	}
	if opts.Keys == "" && opts.Prefix == "" {
		return fmt.Errorf("extract: at least one of --keys or --prefix is required")
	}

	var keys []string
	if opts.Keys != "" {
		for _, k := range strings.Split(opts.Keys, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				keys = append(keys, k)
			}
		}
	}

	exOpts := ExtractOptions{
		Keys:   keys,
		Prefix: opts.Prefix,
		Strip:  opts.Strip,
		DryRun: opts.DryRun,
		Output: opts.Output,
	}

	res, err := Extract(opts.Source, exOpts)
	if err != nil {
		return err
	}

	switch opts.Format {
	case "json":
		return WriteExtractJSON(w, res)
	default:
		WriteExtractText(w, res)
		return nil
	}
}
