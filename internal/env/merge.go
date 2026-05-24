package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// MergeOptions controls how multiple env maps are merged into one file.
type MergeOptions struct {
	// OutputPath is the file to write the merged result to.
	OutputPath string
	// Overwrite controls whether existing keys from later envs replace earlier ones.
	Overwrite bool
	// Header is an optional comment written at the top of the file.
	Header string
	// ExcludeKeys is a set of keys to omit from the output.
	ExcludeKeys map[string]struct{}
}

// DefaultMergeOptions returns sensible defaults.
func DefaultMergeOptions() MergeOptions {
	return MergeOptions{
		Overwrite:   false,
		ExcludeKeys: make(map[string]struct{}),
	}
}

// Merge combines multiple env maps into a single output file.
// Keys appearing in earlier maps take precedence unless Overwrite is true.
// Returns the merged map for inspection.
func Merge(envs []map[string]string, opts MergeOptions) (map[string]string, error) {
	merged := make(map[string]string)

	for _, env := range envs {
		for k, v := range env {
			if _, excluded := opts.ExcludeKeys[k]; excluded {
				continue
			}
			if _, exists := merged[k]; exists && !opts.Overwrite {
				continue
			}
			merged[k] = v
		}
	}

	if opts.OutputPath == "" {
		return merged, nil
	}

	if err := writeMergedFile(merged, opts); err != nil {
		return merged, fmt.Errorf("merge: write %s: %w", opts.OutputPath, err)
	}
	return merged, nil
}

func writeMergedFile(merged map[string]string, opts MergeOptions) error {
	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	if opts.Header != "" {
		for _, line := range strings.Split(opts.Header, "\n") {
			fmt.Fprintf(&sb, "# %s\n", line)
		}
		sb.WriteByte('\n')
	}

	for _, k := range keys {
		v := merged[k]
		if strings.ContainsAny(v, " \t") {
			v = fmt.Sprintf("%q", v)
		}
		fmt.Fprintf(&sb, "%s=%s\n", k, v)
	}

	return os.WriteFile(opts.OutputPath, []byte(sb.String()), 0o644)
}
