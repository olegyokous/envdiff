package env

import (
	"fmt"
	"io"
	"strings"
)

// MaskCmdOptions extends MaskOptions with CLI-level formatting.
type MaskCmdOptions struct {
	MaskOptions
	Format string // "text" or "json"
	Out    io.Writer
}

// RunMask is the CLI entry point for the mask subcommand.
func RunMask(source string, rawKeys string, opts MaskCmdOptions) error {
	if source == "" {
		return fmt.Errorf("mask: source file is required")
	}
	keys := splitKeys(rawKeys)
	if len(keys) == 0 {
		return fmt.Errorf("mask: at least one key must be specified")
	}
	opts.MaskOptions.Keys = keys

	result, err := Mask(source, opts.MaskOptions)
	if err != nil {
		return err
	}

	switch strings.ToLower(opts.Format) {
	case "json":
		return WriteMaskJSON(opts.Out, result)
	default:
		WriteMaskText(opts.Out, result)
		return nil
	}
}

// splitKeys splits a comma-separated key list, trimming whitespace.
func splitKeys(raw string) []string {
	parts := strings.Split(raw, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
