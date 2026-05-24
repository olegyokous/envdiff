package env

import (
	"fmt"
	"os"
	"sort"

	"github.com/your-org/envdiff/internal/parser"
)

// SplitOptions controls how a single .env file is split into multiple output files.
type SplitOptions struct {
	// Prefixes maps a prefix string to an output file path.
	// Keys matching a prefix (case-sensitive) are written to the corresponding file.
	// Keys that match no prefix are written to the Remainder file, if set.
	Prefixes  map[string]string
	Remainder string // path for unmatched keys; ignored if empty
	Strip     bool   // remove the prefix from the key name in the output file
}

// Split reads src and partitions its keys into separate files based on Prefixes.
// It returns the number of keys written to each output path.
func Split(src string, opts SplitOptions) (map[string]int, error) {
	env, err := parser.ParseFile(src)
	if err != nil {
		return nil, fmt.Errorf("split: parse %s: %w", src, err)
	}

	// bucket[outputPath] = ordered list of (key, value)
	type kv struct{ k, v string }
	buckets := map[string][]kv{}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := env[key]
		matched := false
		for prefix, outPath := range opts.Prefixes {
			if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
				outKey := key
				if opts.Strip {
					outKey = key[len(prefix):]
				}
				buckets[outPath] = append(buckets[outPath], kv{outKey, val})
				matched = true
				break
			}
		}
		if !matched && opts.Remainder != "" {
			buckets[opts.Remainder] = append(buckets[opts.Remainder], kv{key, val})
		}
	}

	counts := map[string]int{}
	for outPath, pairs := range buckets {
		if err := writeKVFile(outPath, pairs); err != nil {
			return counts, fmt.Errorf("split: write %s: %w", outPath, err)
		}
		counts[outPath] = len(pairs)
	}
	return counts, nil
}

func writeKVFile(path string, pairs []struct{ k, v string }) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, p := range pairs {
		val := p.v
		if needsQuotes(val) {
			val = `"` + val + `"`
		}
		if _, err := fmt.Fprintf(f, "%s=%s\n", p.k, val); err != nil {
			return err
		}
	}
	return nil
}

func needsQuotes(v string) bool {
	for _, c := range v {
		if c == ' ' || c == '\t' || c == '#' {
			return true
		}
	}
	return false
}
