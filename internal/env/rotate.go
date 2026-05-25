package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// RotateOptions controls the behaviour of the Rotate operation.
type RotateOptions struct {
	// Keys is the list of keys whose values should be replaced.
	Keys []string
	// Values maps key -> new value. If a key is in Keys but not in Values,
	// the new value defaults to an empty string.
	Values map[string]string
	// DryRun skips writing the output file.
	DryRun bool
	// Overwrite controls whether existing keys not in Keys are preserved.
	Overwrite bool
}

// RotateRecord captures what happened to a single key during rotation.
type RotateRecord struct {
	Key      string
	OldValue string
	NewValue string
	Rotated  bool
}

// RotateResult is the outcome of a Rotate call.
type RotateResult struct {
	Source  string
	Records []RotateRecord
}

// Rotate replaces the values of the specified keys in src, writing the result
// to dst unless DryRun is set.
func Rotate(src, dst string, opts RotateOptions) (RotateResult, error) {
	env, err := parser.ParseFile(src)
	if err != nil {
		return RotateResult{}, fmt.Errorf("rotate: parse %s: %w", src, err)
	}

	result := RotateResult{Source: src}
	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	updated := make(map[string]string, len(env))
	for k, v := range env {
		updated[k] = v
	}

	for _, k := range opts.Keys {
		old := env[k]
		newVal := ""
		if opts.Values != nil {
			if nv, ok := opts.Values[k]; ok {
				newVal = nv
			}
		}
		updated[k] = newVal
		result.Records = append(result.Records, RotateRecord{
			Key:      k,
			OldValue: old,
			NewValue: newVal,
			Rotated:  true,
		})
	}

	if opts.DryRun {
		return result, nil
	}

	if err := writeRotateFile(dst, updated); err != nil {
		return result, fmt.Errorf("rotate: write %s: %w", dst, err)
	}
	return result, nil
}

func writeRotateFile(path string, env map[string]string) error {
	var sb strings.Builder
	for k, v := range env {
		if strings.ContainsAny(v, " \t") {
			fmt.Fprintf(&sb, "%s=\"%s\"\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
