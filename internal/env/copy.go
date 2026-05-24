package env

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/your-org/envdiff/internal/parser"
)

// CopyOptions controls how keys are copied between env files.
type CopyOptions struct {
	// Keys is the explicit list of keys to copy. If empty, all keys are copied.
	Keys []string
	// Overwrite replaces existing keys in the destination.
	Overwrite bool
	// DryRun reports what would be written without modifying the file.
	DryRun bool
}

// CopyRecord describes the outcome of copying a single key.
type CopyRecord struct {
	Key     string
	Value   string
	Action  string // "copied", "skipped", "overwritten"
}

// CopyResult holds the full result of a Copy operation.
type CopyResult struct {
	Source      string
	Destination string
	Records     []CopyRecord
}

// Copy reads keys from src and writes them into dst according to opts.
func Copy(src, dst string, opts CopyOptions) (CopyResult, error) {
	srcEnv, err := parser.ParseFile(src)
	if err != nil {
		return CopyResult{}, fmt.Errorf("copy: read source %q: %w", src, err)
	}

	dstEnv, err := parser.ParseFile(dst)
	if err != nil && !os.IsNotExist(err) {
		return CopyResult{}, fmt.Errorf("copy: read destination %q: %w", dst, err)
	}
	if dstEnv == nil {
		dstEnv = map[string]string{}
	}

	wantKeys := opts.Keys
	if len(wantKeys) == 0 {
		for k := range srcEnv {
			wantKeys = append(wantKeys, k)
		}
		sort.Strings(wantKeys)
	}

	result := CopyResult{Source: src, Destination: dst}

	for _, k := range wantKeys {
		v, ok := srcEnv[k]
		if !ok {
			continue
		}
		_, exists := dstEnv[k]
		switch {
		case exists && !opts.Overwrite:
			result.Records = append(result.Records, CopyRecord{Key: k, Value: v, Action: "skipped"})
		case exists && opts.Overwrite:
			dstEnv[k] = v
			result.Records = append(result.Records, CopyRecord{Key: k, Value: v, Action: "overwritten"})
		default:
			dstEnv[k] = v
			result.Records = append(result.Records, CopyRecord{Key: k, Value: v, Action: "copied"})
		}
	}

	if !opts.DryRun {
		if err := writeCopiedFile(dst, dstEnv); err != nil {
			return result, fmt.Errorf("copy: write destination %q: %w", dst, err)
		}
	}
	return result, nil
}

func writeCopiedFile(path string, env map[string]string) error {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := env[k]
		if needsQuotes(v) {
			v = `"` + v + `"`
		}
		sb.WriteString(k + "=" + v + "\n")
	}
	return os.WriteFile(path, []byte(sb.String()), 0644)
}
