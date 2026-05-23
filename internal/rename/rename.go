// Package rename detects keys that have been renamed across environments
// by matching values and flagging keys that appear under different names.
package rename

import (
	"fmt"
	"sort"
)

// Candidate represents a potential key rename between two environments.
type Candidate struct {
	OldKey string
	NewKey string
	Value  string
	EnvA   string
	EnvB   string
}

func (c Candidate) String() string {
	return fmt.Sprintf("%s:%s -> %s:%s (value=%q)", c.EnvA, c.OldKey, c.EnvB, c.NewKey, c.Value)
}

// Detect compares two named env maps and returns candidates where a value
// exists in both envs but under different key names. Only unique values are
// considered to avoid false positives from common values like "true" or "1".
func Detect(nameA string, envA map[string]string, nameB string, envB map[string]string) []Candidate {
	// Build value -> key index for each env, skipping duplicate values.
	indexA := uniqueValueIndex(envA)
	indexB := uniqueValueIndex(envB)

	var candidates []Candidate

	for val, keyA := range indexA {
		keyB, ok := indexB[val]
		if !ok {
			continue
		}
		if keyA == keyB {
			continue // same key, not a rename
		}
		// Only flag if keyA is missing from envB or keyB is missing from envA.
		_, keyAInB := envB[keyA]
		_, keyBInA := envA[keyB]
		if keyAInB || keyBInA {
			continue // both keys exist; ambiguous
		}
		candidates = append(candidates, Candidate{
			OldKey: keyA,
			NewKey: keyB,
			Value:  val,
			EnvA:   nameA,
			EnvB:   nameB,
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].OldKey != candidates[j].OldKey {
			return candidates[i].OldKey < candidates[j].OldKey
		}
		return candidates[i].NewKey < candidates[j].NewKey
	})
	return candidates
}

// uniqueValueIndex builds a map of value -> key, omitting values that appear
// more than once (to avoid false positive rename matches).
func uniqueValueIndex(env map[string]string) map[string]string {
	count := make(map[string]int, len(env))
	for _, v := range env {
		if v != "" {
			count[v]++
		}
	}
	index := make(map[string]string, len(env))
	for k, v := range env {
		if v != "" && count[v] == 1 {
			index[v] = k
		}
	}
	return index
}
