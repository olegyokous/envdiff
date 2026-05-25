package env

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/your/envdiff/internal/parser"
)

// FlattenOptions controls how nested key prefixes are flattened.
type FlattenOptions struct {
	Source    string
	Output    string
	Separator string // e.g. "__" to split APP__DB__HOST into APP.DB.HOST
	Delimiter string // output delimiter, defaults to "_"
	DryRun    bool
	UpperKeys bool
}

// FlattenResult holds the outcome of a flatten operation.
type FlattenResult struct {
	Source   string
	Original map[string]string
	Flat     map[string]string
	Changed  []string
}

// DefaultFlattenOptions returns sensible defaults.
func DefaultFlattenOptions() FlattenOptions {
	return FlattenOptions{
		Separator: "__",
		Delimiter: "_",
	}
}

// Flatten reads a .env file and rewrites keys that contain the separator
// into a flat representation using the delimiter.
func Flatten(opts FlattenOptions) (*FlattenResult, error) {
	if opts.Separator == "" {
		opts.Separator = "__"
	}
	if opts.Delimiter == "" {
		opts.Delimiter = "_"
	}

	original, err := parser.ParseFile(opts.Source)
	if err != nil {
		return nil, fmt.Errorf("flatten: parse %s: %w", opts.Source, err)
	}

	flat := make(map[string]string, len(original))
	var changed []string

	for k, v := range original {
		newKey := strings.ReplaceAll(k, opts.Separator, opts.Delimiter)
		if opts.UpperKeys {
			newKey = strings.ToUpper(newKey)
		}
		flat[newKey] = v
		if newKey != k {
			changed = append(changed, k)
		}
	}

	result := &FlattenResult{
		Source:   opts.Source,
		Original: original,
		Flat:     flat,
		Changed:  changed,
	}

	if !opts.DryRun && opts.Output != "" {
		if err := writeFlatFile(opts.Output, flat); err != nil {
			return result, err
		}
	}

	return result, nil
}

func writeFlatFile(path string, env map[string]string) error {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := env[k]
		if strings.ContainsAny(v, " \t") {
			v = `"` + v + `"`
		}
		sb.WriteString(k + "=" + v + "\n")
	}

	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
