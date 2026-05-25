package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// DedupeOptions controls deduplication behaviour.
type DedupeOptions struct {
	// KeepFirst retains the first occurrence; when false the last wins.
	KeepFirst bool
	// DryRun reports what would change without writing.
	DryRun bool
}

// DedupeResult holds the outcome of a deduplication run.
type DedupeResult struct {
	Source     string
	Duplicates []string // keys that appeared more than once
	Lines      []string // final file lines after dedup
	DryRun     bool
}

// Dedupe reads src, removes duplicate key definitions and writes the result
// back to src (unless DryRun is set).
func Dedupe(src string, opts DedupeOptions) (*DedupeResult, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, fmt.Errorf("dedupe: open %s: %w", src, err)
	}
	defer f.Close()

	type entry struct {
		line string
		key  string
	}

	var entries []entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)
		key := ""
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			if idx := strings.IndexByte(trimmed, '='); idx > 0 {
				key = strings.TrimSpace(trimmed[:idx])
			}
		}
		entries = append(entries, entry{line: raw, key: key})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("dedupe: scan %s: %w", src, err)
	}

	seen := map[string]int{} // key -> count
	for _, e := range entries {
		if e.key != "" {
			seen[e.key]++
		}
	}

	dupeSet := map[string]bool{}
	for k, c := range seen {
		if c > 1 {
			dupeSet[k] = true
		}
	}

	// Decide which occurrence to keep.
	keptCount := map[string]int{}
	var kept []string
	for _, e := range entries {
		if e.key != "" && dupeSet[e.key] {
			keptCount[e.key]++
			keepThis := (opts.KeepFirst && keptCount[e.key] == 1) ||
				(!opts.KeepFirst && keptCount[e.key] == seen[e.key])
			if !keepThis {
				continue
			}
		}
		kept = append(kept, e.line)
	}

	dupes := make([]string, 0, len(dupeSet))
	for k := range dupeSet {
		dupes = append(dupes, k)
	}

	res := &DedupeResult{
		Source:     src,
		Duplicates: dupes,
		Lines:      kept,
		DryRun:     opts.DryRun,
	}

	if !opts.DryRun {
		if err := writeDedupeFile(src, kept); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func writeDedupeFile(path string, lines []string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("dedupe: write %s: %w", path, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, l := range lines {
		fmt.Fprintln(w, l)
	}
	return w.Flush()
}
