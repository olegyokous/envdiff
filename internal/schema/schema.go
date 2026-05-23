// Package schema validates env files against a declared schema,
// ensuring required keys exist and values match expected types.
package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

// FieldType represents the expected value type for a key.
type FieldType string

const (
	TypeString FieldType = "string"
	TypeInt    FieldType = "int"
	TypeBool   FieldType = "bool"
	TypeURL    FieldType = "url"
)

// Field describes a single key in the schema.
type Field struct {
	Key      string    `json:"key"`
	Type     FieldType `json:"type"`
	Required bool      `json:"required"`
	Pattern  string    `json:"pattern,omitempty"`
}

// Schema is a collection of field definitions.
type Schema struct {
	Fields []Field `json:"fields"`
}

// Violation describes a single schema violation.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// Load reads a JSON schema file from disk.
func Load(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("schema: read %q: %w", path, err)
	}
	var s Schema
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("schema: parse %q: %w", path, err)
	}
	return &s, nil
}

// Check validates an env map against the schema and returns all violations.
func (s *Schema) Check(env map[string]string) []Violation {
	var violations []Violation
	for _, f := range s.Fields {
		val, ok := env[f.Key]
		if !ok || val == "" {
			if f.Required {
				violations = append(violations, Violation{Key: f.Key, Message: "required key is missing or empty"})
			}
			continue
		}
		if err := checkType(f.Key, val, f.Type); err != nil {
			violations = append(violations, *err)
		}
		if f.Pattern != "" {
			if matched, _ := regexp.MatchString("^"+f.Pattern+"$", val); !matched {
				violations = append(violations, Violation{Key: f.Key, Message: fmt.Sprintf("value %q does not match pattern %q", val, f.Pattern)})
			}
		}
	}
	return violations
}

var (
	reInt  = regexp.MustCompile(`^-?[0-9]+$`)
	reBool = regexp.MustCompile(`^(true|false|1|0|yes|no)$`)
	reURL  = regexp.MustCompile(`^https?://`)
)

func checkType(key, val string, t FieldType) *Violation {
	switch t {
	case TypeInt:
		if !reInt.MatchString(val) {
			return &Violation{Key: key, Message: fmt.Sprintf("expected int, got %q", val)}
		}
	case TypeBool:
		if !reBool.MatchString(val) {
			return &Violation{Key: key, Message: fmt.Sprintf("expected bool, got %q", val)}
		}
	case TypeURL:
		if !reURL.MatchString(val) {
			return &Violation{Key: key, Message: fmt.Sprintf("expected URL, got %q", val)}
		}
	}
	return nil
}
