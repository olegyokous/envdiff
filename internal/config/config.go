// Package config parses and validates CLI configuration for envdiff.
package config

import (
	"errors"
	"flag"
	"strings"

	"github.com/your-org/envdiff/internal/output"
	"github.com/your-org/envdiff/internal/redact"
)

// Config holds all resolved configuration for a single envdiff run.
type Config struct {
	Files         []string
	Format        output.Format
	StatusFilter  string
	KeyPrefix     string
	KeyPattern    string
	IgnoreFile    string
	BaselineFile  string
	SaveBaseline  bool
	NoSummary     bool
	RedactList    *redact.List
	RedactKeys    []string
}

// Parse reads os.Args via the provided FlagSet and returns a validated Config.
func Parse(fs *flag.FlagSet, args []string) (*Config, error) {
	format := fs.String("format", "text", "output format: text or json")
	status := fs.String("status", "", "filter by status: missing, mismatch, match")
	prefix := fs.String("key-prefix", "", "only include keys with this prefix")
	pattern := fs.String("key-pattern", "", "only include keys matching this regex")
	ignoreFile := fs.String("ignore-file", "", "path to .envignore file")
	baselineFile := fs.String("baseline", "", "path to baseline JSON file")
	saveBaseline := fs.Bool("save-baseline", false, "write current results as new baseline")
	noSummary := fs.Bool("no-summary", false, "suppress summary line")
	redactKeys := fs.String("redact", "", "comma-separated list of additional key patterns to redact")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	files := fs.Args()
	if len(files) < 2 {
		return nil, errors.New("at least two .env files are required")
	}

	fmt, err := output.ParseFormat(*format)
	if err != nil {
		return nil, err
	}

	patterns := redact.DefaultSensitivePatterns
	if *redactKeys != "" {
		for _, k := range strings.Split(*redactKeys, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				patterns = append(patterns, k)
			}
		}
	}

	return &Config{
		Files:        files,
		Format:       fmt,
		StatusFilter: *status,
		KeyPrefix:    *prefix,
		KeyPattern:   *pattern,
		IgnoreFile:   *ignoreFile,
		BaselineFile: *baselineFile,
		SaveBaseline: *saveBaseline,
		NoSummary:    *noSummary,
		RedactList:   redact.NewList(patterns),
		RedactKeys:   patterns,
	}, nil
}
