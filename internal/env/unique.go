package env

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/your/envdiff/internal/parser"
)

// UniqueOptions controls the behaviour of the Unique operation.
type UniqueOptions struct {
	// OnlyIn, when set, restricts output to keys that appear only in that file.
	// Empty string means "keys that are unique across all files" (appear in exactly one).
	OnlyIn string
	DryRun bool
	Output string // destination file; empty means stdout
}

// UniqueRecord describes a single key found to be unique.
type UniqueRecord struct {
	Key      string
	Value    string
	FoundIn  string
}

// UniqueResult is returned by Unique.
type UniqueResult struct {
	Sources []string
	Records []UniqueRecord
}

// Unique finds keys that appear in exactly one of the provided env files.
func Unique(sources []string, opts UniqueOptions) (*UniqueResult, error) {
	if len(sources) < 2 {
		return nil, fmt.Errorf("unique requires at least 2 source files, got %d", len(sources))
	}

	envs := make(map[string]map[string]string, len(sources))
	for _, src := range sources {
		parsed, err := parser.ParseFile(src)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", src, err)
		}
		envs[src] = parsed
	}

	// Count how many files each key appears in.
	keyCount := make(map[string]int)
	keySource := make(map[string]string)
	keyValue := make(map[string]string)
	for src, kv := range envs {
		for k, v := range kv {
			keyCount[k]++
			keySource[k] = src
			keyValue[k] = v
		}
	}

	var records []UniqueRecord
	for k, count := range keyCount {
		if count != 1 {
			continue
		}
		if opts.OnlyIn != "" && keySource[k] != opts.OnlyIn {
			continue
		}
		records = append(records, UniqueRecord{
			Key:     k,
			Value:   keyValue[k],
			FoundIn: keySource[k],
		})
	}
	sort.Slice(records, func(i, j int) bool { return records[i].Key < records[j].Key })

	res := &UniqueResult{Sources: sources, Records: records}

	if !opts.DryRun && opts.Output != "" {
		f, err := os.Create(opts.Output)
		if err != nil {
			return nil, fmt.Errorf("create output: %w", err)
		}
		defer f.Close()
		writeUniqueFile(f, records)
	}

	return res, nil
}

func writeUniqueFile(w io.Writer, records []UniqueRecord) {
	for _, r := range records {
		fmt.Fprintf(w, "%s=%s\n", r.Key, r.Value)
	}
}
