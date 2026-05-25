package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/your-org/envdiff/internal/parser"
)

// TrimOptions controls how unused/empty keys are removed from an env file.
type TrimOptions struct {
	Source    string
	Output    string
	DryRun    bool
	EmptyOnly bool // only remove keys with empty values
}

// TrimRecord describes the action taken on a single key.
type TrimRecord struct {
	Key     string
	Value   string
	Removed bool
}

// TrimResult holds the outcome of a Trim operation.
type TrimResult struct {
	Source  string
	Records []TrimRecord
}

// Trim removes keys from a .env file according to the given options.
// Keys with empty or whitespace-only values are always eligible;
// if EmptyOnly is false, all keys are removed (producing an empty file).
func Trim(opts TrimOptions) (TrimResult, error) {
	env, err := parser.ParseFile(opts.Source)
	if err != nil {
		return TrimResult{}, fmt.Errorf("trim: parse %q: %w", opts.Source, err)
	}

	result := TrimResult{Source: opts.Source}
	kept := map[string]string{}

	for k, v := range env {
		shouldRemove := !opts.EmptyOnly || strings.TrimSpace(v) == ""
		result.Records = append(result.Records, TrimRecord{
			Key:     k,
			Value:   v,
			Removed: shouldRemove,
		})
		if !shouldRemove {
			kept[k] = v
		}
	}

	if opts.DryRun {
		return result, nil
	}

	dest := opts.Output
	if dest == "" {
		dest = opts.Source
	}
	return result, writeTrimFile(dest, kept)
}

func writeTrimFile(path string, env map[string]string) error {
	var sb strings.Builder
	for k, v := range env {
		if strings.ContainsAny(v, " \t") {
			fmt.Fprintf(&sb, "%s=%q\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
