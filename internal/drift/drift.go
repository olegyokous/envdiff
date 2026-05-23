// Package drift detects configuration drift between a live environment
// (expressed as a flat key=value map) and a reference .env file.
package drift

import (
	"fmt"
	"sort"
	"strings"
)

// Status describes the drift state of a single key.
type Status int

const (
	StatusMatch   Status = iota // key exists and value matches
	StatusMissing               // key present in reference but absent in live
	StatusExtra                 // key present in live but absent in reference
	StatusDrifted               // key present in both but values differ
)

func (s Status) String() string {
	switch s {
	case StatusMatch:
		return "match"
	case StatusMissing:
		return "missing"
	case StatusExtra:
		return "extra"
	case StatusDrifted:
		return "drifted"
	default:
		return "unknown"
	}
}

// Entry is a single key's drift result.
type Entry struct {
	Key      string
	Status   Status
	RefValue string // value in the reference file (empty when Extra)
	LiveValue string // value in the live env   (empty when Missing)
}

func (e Entry) String() string {
	return fmt.Sprintf("%s\t%s", e.Key, e.Status)
}

// Compare returns drift entries by comparing reference keys/values against
// the live map. Pass includeExtra=true to also report keys only in live.
func Compare(reference map[string]string, live map[string]string, includeExtra bool) []Entry {
	var entries []Entry

	for k, refVal := range reference {
		liveVal, ok := live[k]
		switch {
		case !ok:
			entries = append(entries, Entry{Key: k, Status: StatusMissing, RefValue: refVal})
		case strings.TrimSpace(liveVal) != strings.TrimSpace(refVal):
			entries = append(entries, Entry{Key: k, Status: StatusDrifted, RefValue: refVal, LiveValue: liveVal})
		default:
			entries = append(entries, Entry{Key: k, Status: StatusMatch, RefValue: refVal, LiveValue: liveVal})
		}
	}

	if includeExtra {
		for k, liveVal := range live {
			if _, found := reference[k]; !found {
				entries = append(entries, Entry{Key: k, Status: StatusExtra, LiveValue: liveVal})
			}
		}
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].Key < entries[j].Key })
	return entries
}

// HasDrift returns true if any entry is not StatusMatch.
func HasDrift(entries []Entry) bool {
	for _, e := range entries {
		if e.Status != StatusMatch {
			return true
		}
	}
	return false
}
