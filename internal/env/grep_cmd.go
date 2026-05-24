// Package env provides utilities for manipulating .env files.
package env

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// GrepCmdOptions configures the grep command execution.
type GrepCmdOptions struct {
	// Sources is the list of .env file paths to search.
	Sources []string
	// Pattern is the search string or regex pattern.
	Pattern string
	// MatchKeys restricts matching to key names only.
	MatchKeys bool
	// MatchValues restricts matching to values only.
	MatchValues bool
	// IgnoreCase enables case-insensitive matching.
	IgnoreCase bool
	// InvertMatch returns lines that do NOT match.
	InvertMatch bool
	// Format is the output format: "text", "json", or "summary".
	Format string
	// Output is the writer for results; defaults to os.Stdout.
	Output io.Writer
}

// DefaultGrepCmdOptions returns sensible defaults for the grep command.
func DefaultGrepCmdOptions() GrepCmdOptions {
	return GrepCmdOptions{
		Format: "text",
		Output: os.Stdout,
	}
}

// RunGrep executes a grep search across one or more .env files and writes
// the results to the configured output in the requested format.
func RunGrep(opts GrepCmdOptions) error {
	if len(opts.Sources) == 0 {
		return fmt.Errorf("grep: at least one source file is required")
	}
	if opts.Pattern == "" {
		return fmt.Errorf("grep: pattern must not be empty")
	}

	w := opts.Output
	if w == nil {
		w = os.Stdout
	}

	grepOpts := DefaultGrepOptions()
	grepOpts.IgnoreCase = opts.IgnoreCase
	grepOpts.InvertMatch = opts.InvertMatch

	if opts.MatchKeys && !opts.MatchValues {
		grepOpts.SearchKeys = true
		grepOpts.SearchValues = false
	} else if opts.MatchValues && !opts.MatchKeys {
		grepOpts.SearchKeys = false
		grepOpts.SearchValues = true
	}
	// If both or neither are set, search both (default behaviour).

	var allMatches []GrepMatch

	for _, src := range opts.Sources {
		matches, err := Grep(src, opts.Pattern, grepOpts)
		if err != nil {
			return fmt.Errorf("grep: reading %s: %w", src, err)
		}
		allMatches = append(allMatches, matches...)
	}

	format := strings.ToLower(strings.TrimSpace(opts.Format))
	switch format {
	case "json":
		return WriteGrepJSON(w, allMatches)
	case "summary":
		return WriteGrepSummary(w, allMatches)
	default:
		WriteGrepText(w, allMatches)
		return nil
	}
}
