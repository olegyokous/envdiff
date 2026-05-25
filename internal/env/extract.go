package env

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/your/envdiff/internal/parser"
)

// ExtractOptions controls the behaviour of Extract.
type ExtractOptions struct {
	Keys     []string // explicit keys to extract; empty means use Prefix
	Prefix   string   // extract all keys with this prefix
	Strip    bool     // strip the prefix from extracted keys
	DryRun   bool
	Output   string // destination file; empty means stdout
}

// ExtractRecord describes the outcome for a single key.
type ExtractRecord struct {
	Key         string
	OriginalKey string // set when Strip was applied
	Value       string
	Found       bool
}

// ExtractResult is returned by Extract.
type ExtractResult struct {
	Source  string
	Records []ExtractRecord
}

// Extract reads src and pulls out the requested keys (or prefix group),
// optionally writing them to a destination file.
func Extract(src string, opts ExtractOptions) (ExtractResult, error) {
	env, err := parser.ParseFile(src)
	if err != nil {
		return ExtractResult{}, fmt.Errorf("extract: parse %s: %w", src, err)
	}

	result := ExtractResult{Source: src}

	if len(opts.Keys) > 0 {
		for _, k := range opts.Keys {
			v, ok := env[k]
			result.Records = append(result.Records, ExtractRecord{
				Key: k, OriginalKey: k, Value: v, Found: ok,
			})
		}
	} else {
		keys := make([]string, 0, len(env))
		for k := range env {
			if opts.Prefix == "" || strings.HasPrefix(k, opts.Prefix) {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		for _, k := range keys {
			out := k
			if opts.Strip && opts.Prefix != "" {
				out = strings.TrimPrefix(k, opts.Prefix)
			}
			result.Records = append(result.Records, ExtractRecord{
				Key: out, OriginalKey: k, Value: env[k], Found: true,
			})
		}
	}

	if opts.DryRun || opts.Output == "" {
		return result, nil
	}
	return result, writeExtractFile(opts.Output, result.Records)
}

func writeExtractFile(path string, records []ExtractRecord) error {
	var sb strings.Builder
	for _, r := range records {
		if !r.Found {
			continue
		}
		v := r.Value
		if needsQuotes(v) {
			v = `"` + v + `"`
		}
		fmt.Fprintf(&sb, "%s=%s\n", r.Key, v)
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
