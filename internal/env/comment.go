package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// CommentOptions controls how comments are added or removed.
type CommentOptions struct {
	Key     string
	Comment string
	Remove  bool
	DryRun  bool
}

// CommentResult describes what happened to a single key's comment.
type CommentResult struct {
	Key     string
	Action  string // "added", "updated", "removed", "not_found"
	Comment string
}

// Comment adds, updates, or removes an inline comment for a key in an env file.
func Comment(src string, opts CommentOptions) (CommentResult, error) {
	f, err := os.Open(src)
	if err != nil {
		return CommentResult{}, fmt.Errorf("open %s: %w", src, err)
	}
	defer f.Close()

	var lines []string
	found := false
	action := "not_found"

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			lines = append(lines, line)
			continue
		}
		eq := strings.IndexByte(line, '=')
		if eq < 0 {
			lines = append(lines, line)
			continue
		}
		key := strings.TrimSpace(line[:eq])
		if key != opts.Key {
			lines = append(lines, line)
			continue
		}
		found = true
		// Strip any existing inline comment from value portion
		valPart := line[eq+1:]
		if idx := strings.Index(valPart, " #"); idx >= 0 {
			valPart = valPart[:idx]
			action = "updated"
		} else {
			action = "added"
		}
		if opts.Remove {
			action = "removed"
			lines = append(lines, key+"="+valPart)
		} else {
			lines = append(lines, key+"="+valPart+" # "+opts.Comment)
		}
	}
	if err := scanner.Err(); err != nil {
		return CommentResult{}, err
	}
	if !found {
		return CommentResult{Key: opts.Key, Action: "not_found"}, nil
	}
	result := CommentResult{Key: opts.Key, Action: action, Comment: opts.Comment}
	if opts.DryRun {
		return result, nil
	}
	if err := writeCommentFile(src, lines); err != nil {
		return CommentResult{}, err
	}
	return result, nil
}

func writeCommentFile(path string, lines []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, l := range lines {
		fmt.Fprintln(w, l)
	}
	return w.Flush()
}
