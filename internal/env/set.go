package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// SetOptions controls the behaviour of Set.
type SetOptions struct {
	File    string
	Key     string
	Value   string
	DryRun  bool
	Create  bool // create the file if it does not exist
}

// SetResult describes what happened during a Set operation.
type SetResult struct {
	File    string
	Key     string
	Value   string
	Prev    string
	Created bool // true when the key did not previously exist
	DryRun  bool
}

// Set writes or updates a single key in an .env file.
func Set(opts SetOptions) (SetResult, error) {
	if opts.Key == "" {
		return SetResult{}, fmt.Errorf("key must not be empty")
	}

	env, err := parser.ParseFile(opts.File)
	if err != nil {
		if !opts.Create {
			return SetResult{}, fmt.Errorf("parse %s: %w", opts.File, err)
		}
		env = map[string]string{}
	}

	prev, existed := env[opts.Key]
	env[opts.Key] = opts.Value

	result := SetResult{
		File:    opts.File,
		Key:     opts.Key,
		Value:   opts.Value,
		Prev:    prev,
		Created: !existed,
		DryRun:  opts.DryRun,
	}

	if opts.DryRun {
		return result, nil
	}

	if err := writeSetFile(opts.File, env); err != nil {
		return SetResult{}, fmt.Errorf("write %s: %w", opts.File, err)
	}
	return result, nil
}

func writeSetFile(path string, env map[string]string) error {
	var sb strings.Builder
	for k, v := range env {
		if needsQuotes(v) {
			fmt.Fprintf(&sb, "%s=\"%s\"\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
