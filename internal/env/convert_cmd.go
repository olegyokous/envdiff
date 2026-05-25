package env

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ConvertCmdOptions configures the convert command.
type ConvertCmdOptions struct {
	Source string
	Dest   string
	Format string
	DryRun bool
	JSON   bool
}

// RunConvert executes the convert command and writes output to w.
func RunConvert(opts ConvertCmdOptions, w io.Writer) error {
	fmt := resolveConvertFormat(opts.Format, opts.Dest)

	r, err := Convert(opts.Source, opts.Dest, fmt, opts.DryRun)
	if err != nil {
		return err
	}

	if opts.JSON {
		return WriteConvertJSON(w, r)
	}
	WriteConvertText(w, r)
	return nil
}

// resolveConvertFormat infers the format from the explicit flag or destination extension.
func resolveConvertFormat(flag, dest string) ConvertFormat {
	if flag != "" {
		switch strings.ToLower(flag) {
		case "export":
			return FormatExport
		case "json":
			return FormatJSON
		case "yaml", "yml":
			return FormatYAML
		default:
			return FormatDotenv
		}
	}
	switch {
	case strings.HasSuffix(dest, ".json"):
		return FormatJSON
	case strings.HasSuffix(dest, ".yaml") || strings.HasSuffix(dest, ".yml"):
		return FormatYAML
	case strings.HasSuffix(dest, ".sh"):
		return FormatExport
	default:
		return FormatDotenv
	}
}

// RunConvertCLI is a thin wrapper for use from main.
func RunConvertCLI(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff convert <source> <dest> [--format=<fmt>] [--dry-run] [--json]")
	}
	opts := ConvertCmdOptions{
		Source: args[0],
		Dest:   args[1],
	}
	for _, a := range args[2:] {
		switch {
		case strings.HasPrefix(a, "--format="):
			opts.Format = strings.TrimPrefix(a, "--format=")
		case a == "--dry-run":
			opts.DryRun = true
		case a == "--json":
			opts.JSON = true
		}
	}
	return RunConvert(opts, os.Stdout)
}
