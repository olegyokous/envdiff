// Package ignore provides functionality for loading and applying
// key ignore lists to diff results, allowing users to suppress
// known or intentional differences from output.
package ignore

import (
	"bufio"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// List holds a set of keys that should be excluded from results.
type List struct {
	keys map[string]struct{}
}

// LoadFile reads an ignore file where each non-blank, non-comment line
// is treated as a key to suppress. Lines starting with '#' are comments.
func LoadFile(path string) (*List, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	l := &List{keys: make(map[string]struct{})}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		l.keys[line] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return l, nil
}

// NewList creates an ignore list from a slice of key strings.
func NewList(keys []string) *List {
	l := &List{keys: make(map[string]struct{}, len(keys))}
	for _, k := range keys {
		l.keys[strings.TrimSpace(k)] = struct{}{}
	}
	return l
}

// Contains reports whether the given key is in the ignore list.
func (l *List) Contains(key string) bool {
	_, ok := l.keys[key]
	return ok
}

// Apply filters out any results whose key appears in the ignore list.
func (l *List) Apply(results []diff.Result) []diff.Result {
	out := make([]diff.Result, 0, len(results))
	for _, r := range results {
		if !l.Contains(r.Key) {
			out = append(out, r)
		}
	}
	return out
}
