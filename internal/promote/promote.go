// Package promote copies approved keys from one environment file to another.
package promote

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Options controls promotion behaviour.
type Options struct {
	// Keys restricts promotion to specific keys; empty means all missing keys.
	Keys []string
	// DryRun reports what would change without writing.
	DryRun bool
	// Overwrite replaces existing keys in the destination.
	Overwrite bool
}

// Result describes the outcome of a single key promotion.
type Result struct {
	Key    string
	Status string // "promoted", "skipped", "dry-run"
}

// Run promotes keys from src env map into the destination file.
// dst is a file path; src is the parsed source environment.
func Run(dstPath string, src map[string]string, dst map[string]string, opts Options) ([]Result, error) {
	keys := candidateKeys(src, dst, opts)

	var results []Result
	updated := copyMap(dst)

	for _, k := range keys {
		v, ok := src[k]
		if !ok {
			continue
		}
		_, exists := dst[k]
		if exists && !opts.Overwrite {
			results = append(results, Result{Key: k, Status: "skipped"})
			continue
		}
		if opts.DryRun {
			results = append(results, Result{Key: k, Status: "dry-run"})
			continue
		}
		updated[k] = v
		results = append(results, Result{Key: k, Status: "promoted"})
	}

	if !opts.DryRun {
		if err := writeEnvFile(dstPath, updated); err != nil {
			return results, fmt.Errorf("promote: write %s: %w", dstPath, err)
		}
	}
	return results, nil
}

func candidateKeys(src, dst map[string]string, opts Options) []string {
	if len(opts.Keys) > 0 {
		return opts.Keys
	}
	var missing []string
	for k := range src {
		if _, ok := dst[k]; !ok {
			missing = append(missing, k)
		}
	}
	sort.Strings(missing)
	return missing
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func writeEnvFile(path string, env map[string]string) error {
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
