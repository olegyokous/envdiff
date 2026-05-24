package env

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// SortOptions controls how a .env file is sorted and written.
type SortOptions struct {
	// Descending sorts keys Z→A instead of A→Z.
	Descending bool
	// DryRun prints the sorted output without writing.
	DryRun bool
	// GroupComments keeps a comment line attached to the key below it.
	GroupComments bool
}

// DefaultSortOptions returns the default sort configuration.
func DefaultSortOptions() SortOptions {
	return SortOptions{}
}

// SortResult holds the outcome of a sort operation.
type SortResult struct {
	Source  string
	KeyCount int
	DryRun  bool
}

// Sort reads src, sorts its key=value lines alphabetically, and writes the
// result back to src (or to stdout when DryRun is true).
func Sort(src string, opts SortOptions) (SortResult, error) {
	f, err := os.Open(src)
	if err != nil {
		return SortResult{}, fmt.Errorf("sort: open %s: %w", src, err)
	}
	defer f.Close()

	type block struct {
		comment string // may be empty
		line    string // key=value
	}

	var blocks []block
	var pending string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			pending = ""
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			if opts.GroupComments {
				pending = raw
			}
			continue
		}
		blocks = append(blocks, block{comment: pending, line: raw})
		pending = ""
	}
	if err := scanner.Err(); err != nil {
		return SortResult{}, fmt.Errorf("sort: scan %s: %w", src, err)
	}

	sort.SliceStable(blocks, func(i, j int) bool {
		ki := strings.SplitN(blocks[i].line, "=", 2)[0]
		kj := strings.SplitN(blocks[j].line, "=", 2)[0]
		if opts.Descending {
			return ki > kj
		}
		return ki < kj
	})

	var sb strings.Builder
	for _, b := range blocks {
		if b.comment != "" {
			sb.WriteString(b.comment + "\n")
		}
		sb.WriteString(b.line + "\n")
	}

	if opts.DryRun {
		fmt.Print(sb.String())
		return SortResult{Source: src, KeyCount: len(blocks), DryRun: true}, nil
	}

	if err := os.WriteFile(src, []byte(sb.String()), 0644); err != nil {
		return SortResult{}, fmt.Errorf("sort: write %s: %w", src, err)
	}
	return SortResult{Source: src, KeyCount: len(blocks)}, nil
}
