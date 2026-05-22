// Package exit provides structured exit codes for CI pipeline integration.
package exit

import (
	"os"

	"github.com/user/envdiff/internal/diff"
)

// Code represents a process exit code.
type Code int

const (
	// CodeOK means all keys matched across all environments.
	CodeOK Code = 0
	// CodeMismatch means one or more keys are missing or have mismatched values.
	CodeMismatch Code = 1
	// CodeError means an unexpected runtime error occurred.
	CodeError Code = 2
)

// FromResults derives the appropriate exit code from a slice of diff results.
// Returns CodeMismatch if any result has a non-match status, otherwise CodeOK.
func FromResults(results []diff.Result) Code {
	for _, r := range results {
		if r.Status != diff.StatusMatch {
			return CodeMismatch
		}
	}
	return CodeOK
}

// Exit terminates the process with the given exit code.
func Exit(code Code) {
	os.Exit(int(code))
}
