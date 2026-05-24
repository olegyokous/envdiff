package env

import (
	"fmt"
	"io"
	"regexp"
	"sort"

	"github.com/your/envdiff/internal/parser"
)

// GrepOptions controls how grep searches env files.
type GrepOptions struct {
	Pattern     string
	SearchKeys  bool
	SearchValues bool
	IgnoreCase  bool
}

// GrepMatch represents a single match found in an env file.
type GrepMatch struct {
	File  string
	Key   string
	Value string
}

// DefaultGrepOptions returns sensible defaults (search both keys and values).
func DefaultGrepOptions() GrepOptions {
	return GrepOptions{
		SearchKeys:   true,
		SearchValues: true,
		IgnoreCase:   false,
	}
}

// Grep searches one or more env files for keys/values matching the given pattern.
func Grep(files []string, opts GrepOptions) ([]GrepMatch, error) {
	if opts.Pattern == "" {
		return nil, fmt.Errorf("grep: pattern must not be empty")
	}

	pattern := opts.Pattern
	if opts.IgnoreCase {
		pattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("grep: invalid pattern %q: %w", opts.Pattern, err)
	}

	var matches []GrepMatch

	for _, file := range files {
		env, err := parser.ParseFile(file)
		if err != nil {
			return nil, fmt.Errorf("grep: reading %s: %w", file, err)
		}

		keys := make([]string, 0, len(env))
		for k := range env {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := env[k]
			keyHit := opts.SearchKeys && re.MatchString(k)
			valHit := opts.SearchValues && re.MatchString(v)
			if keyHit || valHit {
				matches = append(matches, GrepMatch{File: file, Key: k, Value: v})
			}
		}
	}

	return matches, nil
}

// WriteGrepText writes matches in a human-readable format.
func WriteGrepText(w io.Writer, matches []GrepMatch) {
	if len(matches) == 0 {
		fmt.Fprintln(w, "no matches found")
		return
	}
	for _, m := range matches {
		fmt.Fprintf(w, "%s\t%s=%s\n", m.File, m.Key, m.Value)
	}
}
