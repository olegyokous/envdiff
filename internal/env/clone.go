// Package env provides utilities for copying and transforming env maps.
package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// CloneOptions controls how an env map is cloned to a file.
type CloneOptions struct {
	// Keys to exclude from the clone output.
	ExcludeKeys []string
	// If true, values are replaced with empty strings.
	RedactValues bool
	// Optional header comment written at the top of the file.
	Header string
}

// DefaultCloneOptions returns sensible defaults.
func DefaultCloneOptions() CloneOptions {
	return CloneOptions{}
}

// Clone writes a sorted copy of env to path, applying opts.
func Clone(env map[string]string, path string, opts CloneOptions) error {
	exclude := make(map[string]bool, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		exclude[k] = true
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		if !exclude[k] {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var sb strings.Builder
	if opts.Header != "" {
		fmt.Fprintf(&sb, "# %s\n", opts.Header)
	}
	for _, k := range keys {
		v := env[k]
		if opts.RedactValues {
			v = ""
		}
		if strings.ContainsAny(v, " \t") {
			fmt.Fprintf(&sb, "%s=\"%s\"\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}

	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
