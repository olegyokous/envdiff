package env

import (
	"fmt"
	"strings"
)

// TransformOptions controls how key/value transformation is applied.
type TransformOptions struct {
	// KeyPrefix adds a prefix to all keys.
	KeyPrefix string
	// KeySuffix adds a suffix to all keys.
	KeySuffix string
	// UpperKeys converts all keys to uppercase.
	UpperKeys bool
	// LowerKeys converts all keys to lowercase.
	LowerKeys bool
	// ValuePrefix adds a prefix to all values.
	ValuePrefix string
	// ValueSuffix appends a suffix to all values.
	ValueSuffix string
	// RenameMap replaces specific keys; applied after case transforms.
	RenameMap map[string]string
}

// DefaultTransformOptions returns a no-op TransformOptions.
func DefaultTransformOptions() TransformOptions {
	return TransformOptions{}
}

// Apply transforms the provided env map according to the given options.
// Returns a new map; the original is not modified.
func Apply(env map[string]string, opts TransformOptions) (map[string]string, error) {
	if opts.UpperKeys && opts.LowerKeys {
		return nil, fmt.Errorf("transform: UpperKeys and LowerKeys are mutually exclusive")
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		newKey := k
		if opts.UpperKeys {
			newKey = strings.ToUpper(newKey)
		} else if opts.LowerKeys {
			newKey = strings.ToLower(newKey)
		}
		newKey = opts.KeyPrefix + newKey + opts.KeySuffix

		if renamed, ok := opts.RenameMap[newKey]; ok {
			newKey = renamed
		}

		newVal := opts.ValuePrefix + v + opts.ValueSuffix
		out[newKey] = newVal
	}
	return out, nil
}
