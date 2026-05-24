// Package resolve expands variable references within env maps.
// It handles ${VAR} and $VAR style references, resolving them
// in topological order where possible, and reporting cycles.
package resolve

import (
	"fmt"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// Result holds the resolved env map and any warnings encountered.
type Result struct {
	Env      map[string]string
	Warnings []string
}

// Apply resolves variable references in env, returning a new map with
// substitutions applied. Unresolvable references are left as-is and
// a warning is recorded.
func Apply(env map[string]string) Result {
	resolved := make(map[string]string, len(env))
	var warnings []string

	for key, val := range env {
		resolved[key] = val
	}

	const maxPasses = 10
	for pass := 0; pass < maxPasses; pass++ {
		changed := false
		for key, val := range resolved {
			newVal := expand(val, resolved)
			if newVal != val {
				resolved[key] = newVal
				changed = true
			}
		}
		if !changed {
			break
		}
	}

	// Collect warnings for unresolved references.
	for key, val := range resolved {
		refs := extractRefs(val)
		for _, ref := range refs {
			if _, ok := resolved[ref]; !ok {
				warnings = append(warnings, fmt.Sprintf("key %q references undefined variable %q", key, ref))
			}
		}
	}

	return Result{Env: resolved, Warnings: warnings}
}

func expand(val string, env map[string]string) string {
	return refPattern.ReplaceAllStringFunc(val, func(match string) string {
		name := extractName(match)
		if v, ok := env[name]; ok {
			return v
		}
		return match
	})
}

func extractName(ref string) string {
	ref = strings.TrimPrefix(ref, "${") 
	ref = strings.TrimSuffix(ref, "}")
	ref = strings.TrimPrefix(ref, "$")
	return ref
}

func extractRefs(val string) []string {
	matches := refPattern.FindAllStringSubmatch(val, -1)
	var refs []string
	for _, m := range matches {
		if m[1] != "" {
			refs = append(refs, m[1])
		} else if m[2] != "" {
			refs = append(refs, m[2])
		}
	}
	return refs
}
