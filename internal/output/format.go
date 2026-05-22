package output

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the supported output formats.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// ParseFormat converts a string to a Format, returning an error if unknown.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatText, FormatJSON:
		return Format(s), nil
	default:
		return "", fmt.Errorf("unsupported format %q: must be \"text\" or \"json\"", s)
	}
}

// WriterOptions controls where and how output is written.
type WriterOptions struct {
	Format Format
	Out    io.Writer
}

// DefaultOptions returns WriterOptions with text format writing to stdout.
func DefaultOptions() WriterOptions {
	return WriterOptions{
		Format: FormatText,
		Out:    os.Stdout,
	}
}

// Write renders results according to the chosen format.
func Write(results []diff.Result, opts WriterOptions) error {
	w := opts.Out
	if w == nil {
		w = os.Stdout
	}

	switch opts.Format {
	case FormatJSON:
		return diff.WriteJSON(results, w)
	case FormatText:
		return diff.WriteText(results, w)
	default:
		return fmt.Errorf("unknown format: %s", opts.Format)
	}
}

// SupportedFormats returns the list of all recognized format strings.
func SupportedFormats() []Format {
	return []Format{FormatText, FormatJSON}
}
