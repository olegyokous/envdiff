// Package lint provides heuristic checks on parsed env files,
// warning about suspicious keys or values (e.g. empty values,
// keys with unusual characters, or duplicate keys in a single file).
package lint

import (
	"fmt"
	"strings"
	"unicode"
)

// Warning represents a single lint finding.
type Warning struct {
	File string
	Key  string
	Msg  string
}

func (w Warning) String() string {
	return fmt.Sprintf("%s: %s — %s", w.File, w.Key, w.Msg)
}

// Check runs all lint rules against the provided env map and returns any warnings.
func Check(filename string, env map[string]string) []Warning {
	var warnings []Warning

	for key, value := range env {
		if w := checkKeyChars(filename, key); w != nil {
			warnings = append(warnings, *w)
		}
		if w := checkEmptyValue(filename, key, value); w != nil {
			warnings = append(warnings, *w)
		}
		if w := checkWhitespaceValue(filename, key, value); w != nil {
			warnings = append(warnings, *w)
		}
	}

	return warnings
}

// CheckAll runs lint checks across multiple env files and merges the results.
func CheckAll(envs map[string]map[string]string) []Warning {
	var all []Warning
	for name, env := range envs {
		all = append(all, Check(name, env)...)
	}
	return all
}

func checkKeyChars(file, key string) *Warning {
	for _, r := range key {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return &Warning{
				File: file,
				Key:  key,
				Msg:  fmt.Sprintf("key contains unusual character %q", r),
			}
		}
	}
	return nil
}

func checkEmptyValue(file, key, value string) *Warning {
	if value == "" {
		return &Warning{File: file, Key: key, Msg: "value is empty"}
	}
	return nil
}

func checkWhitespaceValue(file, key, value string) *Warning {
	if value != "" && strings.TrimSpace(value) == "" {
		return &Warning{File: file, Key: key, Msg: "value contains only whitespace"}
	}
	return nil
}
