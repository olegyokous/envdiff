package env

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/your-org/envdiff/internal/parser"
)

// IntersectResult holds the outcome of an intersection operation.
type IntersectResult struct {
	Source      string
	Target      string
	Kept        []string
	Dropped     []string
	OutputFile  string
	DryRun      bool
}

// DefaultIntersectOptions returns sensible defaults.
func DefaultIntersectOptions() IntersectOptions {
	return IntersectOptions{
		DryRun: false,
	}
}

// IntersectOptions controls intersection behaviour.
type IntersectOptions struct {
	OutputFile string
	DryRun     bool
}

// Intersect retains only keys that exist in both source and reference files.
// Values are taken from source.
func Intersect(sourcePath, referencePath string, opts IntersectOptions) (IntersectResult, error) {
	src, err := parser.ParseFile(sourcePath)
	if err != nil {
		return IntersectResult{}, fmt.Errorf("intersect: read source: %w", err)
	}

	ref, err := parser.ParseFile(referencePath)
	if err != nil {
		return IntersectResult{}, fmt.Errorf("intersect: read reference: %w", err)
	}

	result := IntersectResult{
		Source:     sourcePath,
		Target:     referencePath,
		OutputFile: opts.OutputFile,
		DryRun:     opts.DryRun,
	}

	kept := make(map[string]string)
	for k, v := range src {
		if _, inRef := ref[k]; inRef {
			kept[k] = v
			result.Kept = append(result.Kept, k)
		} else {
			result.Dropped = append(result.Dropped, k)
		}
	}

	sort.Strings(result.Kept)
	sort.Strings(result.Dropped)

	outPath := opts.OutputFile
	if outPath == "" {
		outPath = sourcePath
	}

	if !opts.DryRun {
		if err := writeIntersectFile(outPath, kept, result.Kept); err != nil {
			return result, fmt.Errorf("intersect: write output: %w", err)
		}
	}

	return result, nil
}

func writeIntersectFile(path string, data map[string]string, keys []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, k := range keys {
		v := data[k]
		if strings.ContainsAny(v, " \t") {
			fmt.Fprintf(f, "%s=\"%s\"\n", k, v)
		} else {
			fmt.Fprintf(f, "%s=%s\n", k, v)
		}
	}
	return nil
}
