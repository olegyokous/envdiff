package env

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/yourorg/envdiff/internal/parser"
)

// RenameOptions controls the rename-key operation.
type RenameOptions struct {
	Source  string
	Dest    string
	OldKey  string
	NewKey  string
	DryRun  bool
	Verbose bool
}

// DefaultRenameOptions returns sensible defaults.
func DefaultRenameOptions() RenameOptions {
	return RenameOptions{}
}

// RenameKey reads Source, renames OldKey to NewKey, and writes result to Dest.
// If DryRun is true the file is not written; changes are printed to w instead.
func RenameKey(opts RenameOptions, w io.Writer) error {
	if opts.OldKey == "" || opts.NewKey == "" {
		return fmt.Errorf("rename: old-key and new-key must not be empty")
	}
	if opts.Source == "" {
		return fmt.Errorf("rename: source file must be specified")
	}

	env, err := parser.ParseFile(opts.Source)
	if err != nil {
		return fmt.Errorf("rename: parse %s: %w", opts.Source, err)
	}

	val, exists := env[opts.OldKey]
	if !exists {
		return fmt.Errorf("rename: key %q not found in %s", opts.OldKey, opts.Source)
	}

	if _, conflict := env[opts.NewKey]; conflict {
		return fmt.Errorf("rename: key %q already exists in %s", opts.NewKey, opts.Source)
	}

	delete(env, opts.OldKey)
	env[opts.NewKey] = val

	if opts.Verbose || opts.DryRun {
		fmt.Fprintf(w, "rename: %s -> %s (value=%q)\n", opts.OldKey, opts.NewKey, val)
	}

	if opts.DryRun {
		return nil
	}

	dest := opts.Dest
	if dest == "" {
		dest = opts.Source
	}
	return writeRenamedFile(dest, env)
}

func writeRenamedFile(path string, env map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		if needsQuotes(v) {
			fmt.Fprintf(f, "%s=%q\n", k, v)
		} else {
			fmt.Fprintf(f, "%s=%s\n", k, v)
		}
	}
	return nil
}
