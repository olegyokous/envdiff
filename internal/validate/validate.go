// Package validate checks that all keys in a reference .env file are present
// in every other environment, and that values conform to optional regex rules.
package validate

import (
	"fmt"
	"regexp"
)

// Rule pairs a key name with a compiled pattern its value must satisfy.
type Rule struct {
	Key     string
	Pattern *regexp.Regexp
}

// Violation describes a single validation failure.
type Violation struct {
	Env     string
	Key     string
	Message string
}

func (v Violation) String() string {
	return fmt.Sprintf("[%s] %s: %s", v.Env, v.Key, v.Message)
}

// Check validates a set of named environments against a list of rules.
// envs maps environment name -> key/value pairs.
// rules is optional; pass nil to skip value-pattern checks.
func Check(envs map[string]map[string]string, rules []Rule) []Violation {
	var violations []Violation

	// Build rule index for O(1) lookup.
	ruleIndex := make(map[string]*regexp.Regexp, len(rules))
	for _, r := range rules {
		ruleIndex[r.Key] = r.Pattern
	}

	for envName, kv := range envs {
		for _, r := range rules {
			val, present := kv[r.Key]
			if !present {
				violations = append(violations, Violation{
					Env:     envName,
					Key:     r.Key,
					Message: "key is missing",
				})
				continue
			}
			if r.Pattern != nil && !r.Pattern.MatchString(val) {
				violations = append(violations, Violation{
					Env:     envName,
					Key:     r.Key,
					Message: fmt.Sprintf("value %q does not match pattern %s", val, r.Pattern),
				})
			}
		}
	}

	return violations
}

// ParseRules converts a map[string]string of key->pattern strings into Rules.
// Returns an error if any pattern fails to compile.
func ParseRules(raw map[string]string) ([]Rule, error) {
	rules := make([]Rule, 0, len(raw))
	for key, pattern := range raw {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern for key %q: %w", key, err)
		}
		rules = append(rules, Rule{Key: key, Pattern: re})
	}
	return rules, nil
}
