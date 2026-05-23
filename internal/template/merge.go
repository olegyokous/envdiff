package template

import (
	"fmt"
	"sort"
)

// MergeResult holds the outcome of merging multiple env maps into a union.
type MergeResult struct {
	// Union contains every key found across all envs. Values come from the
	// first env that defines the key (priority order).
	Union map[string]string
	// OnlyIn maps each key to the list of env names that define it.
	OnlyIn map[string][]string
}

// Merge combines multiple named env maps into a union map.
// envs is a slice of (name, map) pairs represented as a slice of structs.
func Merge(envs []NamedEnv) MergeResult {
	union := make(map[string]string)
	onlyIn := make(map[string][]string)

	for _, ne := range envs {
		for k, v := range ne.Env {
			if _, exists := union[k]; !exists {
				union[k] = v
			}
			onlyIn[k] = append(onlyIn[k], ne.Name)
		}
	}

	// Sort the name lists for deterministic output.
	for k := range onlyIn {
		sort.Strings(onlyIn[k])
	}

	return MergeResult{Union: union, OnlyIn: onlyIn}
}

// NamedEnv pairs an environment label with its key/value map.
type NamedEnv struct {
	Name string
	Env  map[string]string
}

// UniversalKeys returns keys that appear in every one of the provided envs.
func UniversalKeys(envs []NamedEnv) []string {
	if len(envs) == 0 {
		return nil
	}
	counts := make(map[string]int)
	for _, ne := range envs {
		for k := range ne.Env {
			counts[k]++
		}
	}
	var keys []string
	for k, c := range counts {
		if c == len(envs) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}

// FormatOnlyIn returns a human-readable string for keys missing from some envs.
func FormatOnlyIn(key string, names []string) string {
	return fmt.Sprintf("%s: present in [%v]", key, names)
}
