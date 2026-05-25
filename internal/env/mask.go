package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/your/envdiff/internal/parser"
)

// MaskOptions controls how masking is applied.
type MaskOptions struct {
	// Keys whose values will be masked.
	Keys []string
	// MaskChar is the character used for masking (default "*").
	MaskChar string
	// MaskLen is the fixed length of the mask (0 = full replacement).
	MaskLen int
	// DryRun prevents writing to disk.
	DryRun bool
	// Output path; empty means overwrite Source.
	Output string
}

// MaskRecord describes what happened to a single key.
type MaskRecord struct {
	Key    string
	Masked bool
}

// MaskResult is the outcome of a Mask operation.
type MaskResult struct {
	Source  string
	Records []MaskRecord
}

// Mask replaces the values of the specified keys in the env file with a mask string.
func Mask(source string, opts MaskOptions) (MaskResult, error) {
	if opts.MaskChar == "" {
		opts.MaskChar = "*"
	}

	env, err := parser.ParseFile(source)
	if err != nil {
		return MaskResult{}, fmt.Errorf("mask: parse %s: %w", source, err)
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	var records []MaskRecord
	masked := make(map[string]string, len(env))
	for k, v := range env {
		if _, hit := keySet[k]; hit {
			l := opts.MaskLen
			if l <= 0 {
				l = len(v)
				if l == 0 {
					l = 8
				}
			}
			masked[k] = strings.Repeat(opts.MaskChar, l)
			records = append(records, MaskRecord{Key: k, Masked: true})
		} else {
			masked[k] = v
			records = append(records, MaskRecord{Key: k, Masked: false})
		}
	}

	result := MaskResult{Source: source, Records: records}

	if opts.DryRun {
		return result, nil
	}

	dest := opts.Output
	if dest == "" {
		dest = source
	}
	if err := writeMaskFile(dest, masked); err != nil {
		return result, err
	}
	return result, nil
}

func writeMaskFile(path string, env map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("mask: create %s: %w", path, err)
	}
	defer f.Close()
	for k, v := range env {
		if _, err := fmt.Fprintf(f, "%s=%s\n", k, v); err != nil {
			return err
		}
	}
	return nil
}
