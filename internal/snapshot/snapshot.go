// Package snapshot captures and compares point-in-time states of env files.
package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/user/envdiff/internal/parser"
)

// Entry holds a single snapshot record.
type Entry struct {
	TakenAt time.Time            `json:"taken_at"`
	Label   string               `json:"label"`
	Files   map[string]string    `json:"files"` // filename -> checksum
	Envs    map[string]EnvState  `json:"envs"`
}

// EnvState holds the parsed key/value pairs and a content checksum.
type EnvState struct {
	Checksum string            `json:"checksum"`
	Keys     map[string]string `json:"keys"`
}

// Diff describes keys that changed between two snapshots for one env file.
type Diff struct {
	File    string
	Added   []string
	Removed []string
	Changed []string
}

// Take parses each file and returns a new Entry.
func Take(label string, files []string) (*Entry, error) {
	entry := &Entry{
		TakenAt: time.Now().UTC(),
		Label:   label,
		Files:   make(map[string]string),
		Envs:    make(map[string]EnvState),
	}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("snapshot: read %s: %w", f, err)
		}
		sum := sha256.Sum256(data)
		entry.Files[f] = hex.EncodeToString(sum[:])

		kvs, err := parser.ParseFile(f)
		if err != nil {
			return nil, fmt.Errorf("snapshot: parse %s: %w", f, err)
		}
		entry.Envs[f] = EnvState{Checksum: entry.Files[f], Keys: kvs}
	}
	return entry, nil
}

// Save writes an Entry to path as JSON.
func Save(path string, e *Entry) error {
	b, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// Load reads an Entry from a JSON file.
func Load(path string) (*Entry, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: load %s: %w", path, err)
	}
	var e Entry
	if err := json.Unmarshal(b, &e); err != nil {
		return nil, fmt.Errorf("snapshot: decode %s: %w", path, err)
	}
	return &e, nil
}

// Compare returns per-file diffs between two snapshots.
func Compare(old, new *Entry) []Diff {
	var diffs []Diff
	for file, newState := range new.Envs {
		oldState, exists := old.Envs[file]
		if !exists {
			keys := sortedKeys(newState.Keys)
			diffs = append(diffs, Diff{File: file, Added: keys})
			continue
		}
		d := Diff{File: file}
		for k, v := range newState.Keys {
			if ov, ok := oldState.Keys[k]; !ok {
				d.Added = append(d.Added, k)
			} else if ov != v {
				d.Changed = append(d.Changed, k)
			}
		}
		for k := range oldState.Keys {
			if _, ok := newState.Keys[k]; !ok {
				d.Removed = append(d.Removed, k)
			}
		}
		sort.Strings(d.Added)
		sort.Strings(d.Removed)
		sort.Strings(d.Changed)
		diffs = append(diffs, d)
	}
	return diffs
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
