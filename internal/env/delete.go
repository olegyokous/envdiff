package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// DeleteOptions controls the behaviour of Delete.
type DeleteOptions struct {
	DryRun bool
	Force  bool // suppress error when key does not exist
}

// DeleteRecord captures what happened to a single key.
type DeleteRecord struct {
	Key     string
	Deleted bool
	Missing bool // key was not present in the file
}

// DeleteResult is the outcome of a Delete operation.
type DeleteResult struct {
	Source  string
	Records []DeleteRecord
	DryRun  bool
}

// Delete removes one or more keys from a .env file.
func Delete(src string, keys []string, opts DeleteOptions) (DeleteResult, error) {
	env, err := parser.ParseFile(src)
	if err != nil {
		return DeleteResult{}, fmt.Errorf("delete: parse %q: %w", src, err)
	}

	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[k] = true
	}

	result := DeleteResult{Source: src, DryRun: opts.DryRun}

	for _, k := range keys {
		if _, ok := env[k]; ok {
			result.Records = append(result.Records, DeleteRecord{Key: k, Deleted: true})
		} else {
			if !opts.Force {
				return DeleteResult{}, fmt.Errorf("delete: key %q not found in %s", k, src)
			}
			result.Records = append(result.Records, DeleteRecord{Key: k, Missing: true})
		}
	}

	if opts.DryRun {
		return result, nil
	}

	// Rebuild the file, preserving order and comments.
	raw, err := os.ReadFile(src)
	if err != nil {
		return DeleteResult{}, fmt.Errorf("delete: read %q: %w", src, err)
	}

	var out strings.Builder
	for _, line := range strings.Split(string(raw), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			out.WriteString(line + "\n")
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) == 2 && keySet[strings.TrimSpace(parts[0])] {
			continue // drop this line
		}
		out.WriteString(line + "\n")
	}

	if err := writeDeleteFile(src, out.String()); err != nil {
		return DeleteResult{}, err
	}
	return result, nil
}

func writeDeleteFile(path, content string) error {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("delete: write %q: %w", path, err)
	}
	return nil
}
