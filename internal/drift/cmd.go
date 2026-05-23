package drift

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/output"
	"github.com/user/envdiff/internal/parser"
)

// Options configures a drift check run.
type Options struct {
	ReferenceFile string
	LiveFile      string
	IncludeExtra  bool
	Format        string // "text" or "json"
	FailOnDrift   bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Format:      "text",
		FailOnDrift: true,
	}
}

// Run executes a drift comparison and writes results to stdout.
// Returns a non-zero exit code if drift is detected and FailOnDrift is set.
func Run(opts Options) int {
	ref, err := parser.ParseFile(opts.ReferenceFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "envdiff drift: reading reference file: %v\n", err)
		return 2
	}

	live, err := parser.ParseFile(opts.LiveFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "envdiff drift: reading live file: %v\n", err)
		return 2
	}

	entries := Compare(ref, live, opts.IncludeExtra)

	fmt := output.ParseFormat(opts.Format)
	switch fmt {
	case output.FormatJSON:
		if err := WriteJSON(os.Stdout, entries); err != nil {
			fmt.Fprintf(os.Stderr, "envdiff drift: writing JSON: %v\n", err)
			return 2
		}
	default:
		WriteText(os.Stdout, entries)
	}

	if opts.FailOnDrift && HasDrift(entries) {
		return 1
	}
	return 0
}
