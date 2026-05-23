// Package redact provides utilities for masking sensitive values in diff output.
package redact

import (
	"strings"

	"github.com/your-org/envdiff/internal/diff"
)

// DefaultSensitivePatterns contains common substrings that indicate a key holds sensitive data.
var DefaultSensitivePatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
}

// Mask is the string used to replace sensitive values.
const Mask = "***REDACTED***"

// List holds a set of key patterns considered sensitive.
type List struct {
	patterns []string
}

// NewList creates a List from the provided patterns (uppercased for comparison).
func NewList(patterns []string) *List {
	norm := make([]string, len(patterns))
	for i, p := range patterns {
		norm[i] = strings.ToUpper(p)
	}
	return &List{patterns: norm}
}

// IsSensitive reports whether the key matches any sensitive pattern.
func (l *List) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range l.patterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// Apply returns a new slice of Results with sensitive values replaced by Mask.
func (l *List) Apply(results []diff.Result) []diff.Result {
	out := make([]diff.Result, len(results))
	for i, r := range results {
		if l.IsSensitive(r.Key) {
			masked := make(map[string]string, len(r.Values))
			for env, val := range r.Values {
				if val != "" {
					masked[env] = Mask
				} else {
					masked[env] = val
				}
			}
			r.Values = masked
		}
		out[i] = r
	}
	return out
}
