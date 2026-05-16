package filter

import (
	"regexp"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Options controls which results are included in filtered output.
type Options struct {
	// StatusFilter limits results to a specific status ("match", "missing", "mismatch").
	// Empty string means no filter.
	StatusFilter string

	// KeyPrefix limits results to keys that start with the given prefix.
	KeyPrefix string

	// KeyPattern limits results to keys matching the given regex pattern.
	KeyPattern string
}

// Apply returns a filtered copy of results based on the provided Options.
func Apply(results []diff.Result, opts Options) ([]diff.Result, error) {
	var pattern *regexp.Regexp
	if opts.KeyPattern != "" {
		var err error
		pattern, err = regexp.Compile(opts.KeyPattern)
		if err != nil {
			return nil, err
		}
	}

	var out []diff.Result
	for _, r := range results {
		if opts.StatusFilter != "" && !strings.EqualFold(r.Status, opts.StatusFilter) {
			continue
		}
		if opts.KeyPrefix != "" && !strings.HasPrefix(r.Key, opts.KeyPrefix) {
			continue
		}
		if pattern != nil && !pattern.MatchString(r.Key) {
			continue
		}
		out = append(out, r)
	}
	return out, nil
}
