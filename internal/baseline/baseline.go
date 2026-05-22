// Package baseline provides functionality for saving and loading a reference
// .env snapshot to compare future runs against, enabling drift detection.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved baseline of environment key/value pairs.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Source    string            `json:"source"`
	Env       map[string]string `json:"env"`
}

// Save writes a snapshot of env to the given file path as JSON.
func Save(path string, source string, env map[string]string) error {
	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		Source:    source,
		Env:       env,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write %s: %w", path, err)
	}

	return nil
}

// Load reads a previously saved snapshot from path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("baseline: file not found: %s", path)
		}
		return nil, fmt.Errorf("baseline: read %s: %w", path, err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("baseline: parse %s: %w", path, err)
	}

	return &snap, nil
}

// Diff returns keys present in current but missing from snap, and keys whose
// values differ between snap and current.
func Diff(snap *Snapshot, current map[string]string) (missing []string, changed []string) {
	for k := range current {
		snapVal, ok := snap.Env[k]
		if !ok {
			missing = append(missing, k)
			continue
		}
		if snapVal != current[k] {
			changed = append(changed, k)
		}
	}
	return missing, changed
}
