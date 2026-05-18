package config

import (
	"errors"
	"flag"
	"strings"
)

// Options holds the parsed CLI configuration for an envdiff run.
type Options struct {
	Files      []string
	Format     string
	StatusFilter string
	KeyPrefix  string
	KeyPattern string
	NoColor    bool
	ExitCode   bool
}

// Parse reads command-line arguments and returns a populated Options struct.
// It returns an error if required arguments are missing or invalid.
func Parse(args []string) (*Options, error) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)

	format := fs.String("format", "text", "Output format: text or json")
	statusFilter := fs.String("status", "", "Filter by status: ok, missing, mismatch")
	keyPrefix := fs.String("prefix", "", "Only include keys with this prefix")
	keyPattern := fs.String("pattern", "", "Only include keys matching this regex pattern")
	noColor := fs.Bool("no-color", false, "Disable colored output")
	exitCode := fs.Bool("exit-code", false, "Exit with non-zero code if differences found")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	files := fs.Args()
	if len(files) < 2 {
		return nil, errors.New("at least two .env files must be provided")
	}

	*format = strings.ToLower(strings.TrimSpace(*format))
	if *format != "text" && *format != "json" {
		return nil, errors.New("invalid format: must be \"text\" or \"json\"")
	}

	validStatuses := map[string]bool{"ok": true, "missing": true, "mismatch": true, "": true}
	if !validStatuses[strings.ToLower(*statusFilter)] {
		return nil, errors.New("invalid status filter: must be ok, missing, or mismatch")
	}

	return &Options{
		Files:        files,
		Format:       *format,
		StatusFilter: strings.ToLower(*statusFilter),
		KeyPrefix:    *keyPrefix,
		KeyPattern:   *keyPattern,
		NoColor:      *noColor,
		ExitCode:     *exitCode,
	}, nil
}
