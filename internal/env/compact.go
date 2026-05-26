package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// CompactResult holds the outcome of a compact operation.
type CompactResult struct {
	Source        string
	Output        string
	RemovedCount  int
	KeptCount     int
	DryRun        bool
	Records       []CompactRecord
}

// CompactRecord describes what happened to a single line.
type CompactRecord struct {
	Key     string
	Removed bool
	Reason  string
}

// Compact removes comment lines and blank lines from an env file,
// producing a minimal key=value only file.
func Compact(src, dst string, dryRun bool) (CompactResult, error) {
	env, err := parser.ParseFile(src)
	if err != nil {
		return CompactResult{}, fmt.Errorf("compact: parse %s: %w", src, err)
	}

	// Re-read raw lines to count what we remove.
	f, err := os.Open(src)
	if err != nil {
		return CompactResult{}, fmt.Errorf("compact: open %s: %w", src, err)
	}
	defer f.Close()

	var rawLines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		rawLines = append(rawLines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return CompactResult{}, fmt.Errorf("compact: read %s: %w", src, err)
	}

	var records []CompactRecord
	removed := 0
	for _, line := range rawLines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			records = append(records, CompactRecord{Key: "", Removed: true, Reason: "blank line"})
			removed++
		} else if strings.HasPrefix(trimmed, "#") {
			records = append(records, CompactRecord{Key: trimmed, Removed: true, Reason: "comment"})
			removed++
		}
	}

	result := CompactResult{
		Source:       src,
		Output:       dst,
		RemovedCount: removed,
		KeptCount:    len(env),
		DryRun:       dryRun,
		Records:      records,
	}

	if !dryRun {
		if err := writeCompactFile(dst, env); err != nil {
			return result, err
		}
	}
	return result, nil
}

func writeCompactFile(path string, env map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("compact: create %s: %w", path, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for k, v := range env {
		if strings.ContainsAny(v, " \t") {
			fmt.Fprintf(w, "%s=%q\n", k, v)
		} else {
			fmt.Fprintf(w, "%s=%s\n", k, v)
		}
	}
	return w.Flush()
}
