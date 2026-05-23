// Package template provides the generate sub-command logic.
package template

import (
	"fmt"
	"io"
	"os"
)

// GenerateCmd holds the parsed arguments for the generate command.
type GenerateCmd struct {
	// InputFiles are the source .env files to merge.
	InputFiles []string
	// OutputPath is where the .env.example will be written ("-" for stdout).
	OutputPath string
	// Redact controls whether values are stripped.
	Redact bool
	// Header is the optional file header comment.
	Header string
}

// Run executes the generate command using the provided parser function.
// parseFile must match the signature of parser.ParseFile.
func Run(cmd GenerateCmd, parseFile func(string) (map[string]string, error)) error {
	if len(cmd.InputFiles) == 0 {
		return fmt.Errorf("template: at least one input file required")
	}

	var envs []NamedEnv
	for _, f := range cmd.InputFiles {
		env, err := parseFile(f)
		if err != nil {
			return fmt.Errorf("template: parse %s: %w", f, err)
		}
		envs = append(envs, NamedEnv{Name: f, Env: env})
	}

	result := Merge(envs)

	opts := Options{
		Redact: cmd.Redact,
		Header: cmd.Header,
	}

	var w io.Writer
	if cmd.OutputPath == "-" || cmd.OutputPath == "" {
		w = os.Stdout
		return Generate(w, result.Union, opts)
	}
	return WriteFile(cmd.OutputPath, result.Union, opts)
}
