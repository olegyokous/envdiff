// Package watch monitors .env files for changes and re-runs comparisons.
package watch

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Options configures the watch behavior.
type Options struct {
	Interval time.Duration
	Out      io.Writer
}

// DefaultOptions returns sensible defaults for watching.
func DefaultOptions() Options {
	return Options{
		Interval: 2 * time.Second,
		Out:      os.Stdout,
	}
}

// FileState holds the last-known modification time of a file.
type FileState struct {
	Path    string
	ModTime time.Time
}

// Snapshot captures the current mod times for a set of files.
func Snapshot(paths []string) ([]FileState, error) {
	states := make([]FileState, 0, len(paths))
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			return nil, fmt.Errorf("watch: stat %s: %w", p, err)
		}
		states = append(states, FileState{Path: p, ModTime: info.ModTime()})
	}
	return states, nil
}

// Changed returns true if any file in next has a newer mod time than prev.
func Changed(prev, next []FileState) bool {
	index := make(map[string]time.Time, len(prev))
	for _, s := range prev {
		index[s.Path] = s.ModTime
	}
	for _, s := range next {
		if t, ok := index[s.Path]; !ok || s.ModTime.After(t) {
			return true
		}
	}
	return false
}

// Watch polls the given paths at opts.Interval and calls onChange whenever
// a change is detected. It blocks until the done channel is closed.
func Watch(paths []string, opts Options, onChange func(), done <-chan struct{}) error {
	prev, err := Snapshot(paths)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return nil
		case <-ticker.C:
			next, err := Snapshot(paths)
			if err != nil {
				fmt.Fprintf(opts.Out, "watch: error: %v\n", err)
				continue
			}
			if Changed(prev, next) {
				fmt.Fprintf(opts.Out, "watch: change detected, re-running...\n")
				onChange()
				prev = next
			}
		}
	}
}
