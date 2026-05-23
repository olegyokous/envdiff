// Package audit provides the CLI command handler for audit log operations.
package audit

import (
	"fmt"
	"io"
)

// Options controls the audit command behaviour.
type Options struct {
	LogPath string
	Format  string // "text" or "json"
	Last    bool   // show only the most recent entry
}

// Run executes the audit sub-command, writing output to w.
func Run(w io.Writer, opts Options) error {
	if opts.LogPath == "" {
		return fmt.Errorf("audit: log path must not be empty")
	}

	entries, err := Load(opts.LogPath)
	if err != nil {
		return fmt.Errorf("audit: loading log: %w", err)
	}

	if opts.Last {
		switch opts.Format {
		case "json":
			if len(entries) > 0 {
				return WriteJSON(w, entries[len(entries)-1:])
			}
			WriteLast(w, entries)
			return nil
		default:
			WriteLast(w, entries)
			return nil
		}
	}

	switch opts.Format {
	case "json":
		return WriteJSON(w, entries)
	default:
		WriteText(w, entries)
		return nil
	}
}
