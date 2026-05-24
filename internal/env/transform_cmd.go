package env

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/parser"
)

// TransformCmdOptions bundles CLI-level parameters for the transform command.
type TransformCmdOptions struct {
	Input      string
	Output     string
	Transform  TransformOptions
	DryRun     bool
}

// RunTransform reads an env file, applies transformations, and writes the result.
func RunTransform(opts TransformCmdOptions, stderr io.Writer) error {
	env, err := parser.ParseFile(opts.Input)
	if err != nil {
		return fmt.Errorf("transform: read %s: %w", opts.Input, err)
	}

	transformed, err := Apply(env, opts.Transform)
	if err != nil {
		return fmt.Errorf("transform: %w", err)
	}

	if opts.DryRun {
		keys := make([]string, 0, len(transformed))
		for k := range transformed {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		fmt.Fprintln(stderr, "# dry-run: transformed output")
		for _, k := range keys {
			fmt.Fprintf(stderr, "%s=%s\n", k, transformed[k])
		}
		return nil
	}

	if err := writeTransformedFile(opts.Output, transformed); err != nil {
		return fmt.Errorf("transform: write %s: %w", opts.Output, err)
	}
	return nil
}

func writeTransformedFile(path string, env map[string]string) error {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	lines := make([]string, 0, len(keys))
	for _, k := range keys {
		lines = append(lines, fmt.Sprintf("%s=%s", k, env[k]))
	}
	return writeKVFile(path, lines)
}
