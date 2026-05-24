package env

import (
	"fmt"
	"io"
	"os"
)

// SetCmdOptions is the CLI-level configuration for the set command.
type SetCmdOptions struct {
	File   string
	Key    string
	Value  string
	DryRun bool
	Create bool
	Format string // "text" | "json"
	Out    io.Writer
}

// RunSet is the entry point called by the CLI for the `set` sub-command.
func RunSet(opts SetCmdOptions) error {
	if opts.Out == nil {
		opts.Out = os.Stdout
	}

	result, err := Set(SetOptions{
		File:   opts.File,
		Key:    opts.Key,
		Value:  opts.Value,
		DryRun: opts.DryRun,
		Create: opts.Create,
	})
	if err != nil {
		return err
	}

	switch opts.Format {
	case "json":
		if err := WriteSetJSON(opts.Out, result); err != nil {
			return fmt.Errorf("encode json: %w", err)
		}
	default:
		WriteSetText(opts.Out, result)
	}
	return nil
}
