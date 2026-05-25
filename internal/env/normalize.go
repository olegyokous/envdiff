package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/jakeloud/envdiff/internal/parser"
)

// NormalizeOptions controls how keys and values are normalized.
type NormalizeOptions struct {
	Source      string
	Output      string
	UpperKeys   bool
	LowerKeys   bool
	TrimValues  bool
	QuoteAll    bool
	DryRun      bool
}

// NormalizeResult holds the outcome of a normalize operation.
type NormalizeResult struct {
	Source  string
	Records []NormalizeRecord
}

// NormalizeRecord describes a single key transformation.
type NormalizeRecord struct {
	Key      string
	OldKey   string
	OldValue string
	NewValue string
	Changed  bool
}

// Normalize reads a .env file and normalizes keys/values according to opts.
func Normalize(opts NormalizeOptions) (NormalizeResult, error) {
	if opts.UpperKeys && opts.LowerKeys {
		return NormalizeResult{}, fmt.Errorf("upper-keys and lower-keys are mutually exclusive")
	}

	env, err := parser.ParseFile(opts.Source)
	if err != nil {
		return NormalizeResult{}, fmt.Errorf("parse %s: %w", opts.Source, err)
	}

	result := NormalizeResult{Source: opts.Source}
	normalized := make(map[string]string, len(env))

	for k, v := range env {
		newKey := k
		if opts.UpperKeys {
			newKey = strings.ToUpper(k)
		} else if opts.LowerKeys {
			newKey = strings.ToLower(k)
		}

		newVal := v
		if opts.TrimValues {
			newVal = strings.TrimSpace(v)
		}
		if opts.QuoteAll && !strings.HasPrefix(newVal, `"`) {
			newVal = `"` + newVal + `"`
		}

		changed := newKey != k || newVal != v
		result.Records = append(result.Records, NormalizeRecord{
			Key:      newKey,
			OldKey:   k,
			OldValue: v,
			NewValue: newVal,
			Changed:  changed,
		})
		normalized[newKey] = newVal
	}

	if !opts.DryRun {
		dest := opts.Output
		if dest == "" {
			dest = opts.Source
		}
		if err := writeNormalizeFile(dest, normalized); err != nil {
			return result, fmt.Errorf("write %s: %w", dest, err)
		}
	}
	return result, nil
}

func writeNormalizeFile(path string, env map[string]string) error {
	var sb strings.Builder
	for k, v := range env {
		fmt.Fprintf(&sb, "%s=%s\n", k, v)
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
