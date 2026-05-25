package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// ConvertFormat represents a target serialization format.
type ConvertFormat string

const (
	FormatDotenv ConvertFormat = "dotenv"
	FormatExport ConvertFormat = "export"
	FormatJSON   ConvertFormat = "json"
	FormatYAML   ConvertFormat = "yaml"
)

// ConvertResult holds the outcome of a conversion operation.
type ConvertResult struct {
	Source  string
	Dest    string
	Format  ConvertFormat
	Keys    int
	DryRun  bool
	Output  string
}

// Convert reads a .env file and writes it in the target format.
func Convert(src string, dest string, format ConvertFormat, dryRun bool) (*ConvertResult, error) {
	env, err := parser.ParseFile(src)
	if err != nil {
		return nil, fmt.Errorf("convert: parse %s: %w", src, err)
	}

	var sb strings.Builder
	switch format {
	case FormatExport:
		for k, v := range env {
			fmt.Fprintf(&sb, "export %s=%q\n", k, v)
		}
	case FormatJSON:
		sb.WriteString("{\n")
		i := 0
		for k, v := range env {
			comma := ","
			if i == len(env)-1 {
				comma = ""
			}
			fmt.Fprintf(&sb, "  %q: %q%s\n", k, v, comma)
			i++
		}
		sb.WriteString("}\n")
	case FormatYAML:
		for k, v := range env {
			fmt.Fprintf(&sb, "%s: %q\n", k, v)
		}
	default: // dotenv
		for k, v := range env {
			if strings.ContainsAny(v, " \t") {
				fmt.Fprintf(&sb, "%s=%q\n", k, v)
			} else {
				fmt.Fprintf(&sb, "%s=%s\n", k, v)
			}
		}
	}

	result := &ConvertResult{
		Source: src,
		Dest:   dest,
		Format: format,
		Keys:   len(env),
		DryRun: dryRun,
		Output: sb.String(),
	}

	if !dryRun {
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return nil, fmt.Errorf("convert: mkdir: %w", err)
		}
		if err := os.WriteFile(dest, []byte(sb.String()), 0o644); err != nil {
			return nil, fmt.Errorf("convert: write %s: %w", dest, err)
		}
	}

	return result, nil
}
