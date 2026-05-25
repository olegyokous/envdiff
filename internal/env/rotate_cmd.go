package env

import (
	"fmt"
	"io"
	"strings"
)

// RotateCmdOptions are the CLI-level options for the rotate command.
type RotateCmdOptions struct {
	Source  string
	Dest    string
	// Assignments is a slice of "KEY=VALUE" strings.
	Assignments []string
	DryRun      bool
	Format      string // "text" or "json"
}

// RunRotate is the entry point called by the CLI rotate sub-command.
func RunRotate(opts RotateCmdOptions, stdout io.Writer) error {
	if opts.Source == "" {
		return fmt.Errorf("rotate: source file is required")
	}
	dst := opts.Dest
	if dst == "" {
		dst = opts.Source
	}

	keys := make([]string, 0, len(opts.Assignments))
	values := make(map[string]string, len(opts.Assignments))
	for _, a := range opts.Assignments {
		parts := strings.SplitN(a, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("rotate: invalid assignment %q (expected KEY=VALUE)", a)
		}
		keys = append(keys, parts[0])
		values[parts[0]] = parts[1]
	}

	if len(keys) == 0 {
		return fmt.Errorf("rotate: at least one KEY=VALUE assignment is required")
	}

	result, err := Rotate(opts.Source, dst, RotateOptions{
		Keys:   keys,
		Values: values,
		DryRun: opts.DryRun,
	})
	if err != nil {
		return err
	}

	switch opts.Format {
	case "json":
		return WriteRotateJSON(stdout, result)
	default:
		WriteRotateText(stdout, result)
		return nil
	}
}
