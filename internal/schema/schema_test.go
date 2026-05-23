package schema_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/schema"
)

func writeSchema(t *testing.T, s schema.Schema) string {
	t.Helper()
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(t.TempDir(), "schema.json")
	if err := os.WriteFile(p, data, 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_ValidFile(t *testing.T) {
	s := schema.Schema{Fields: []schema.Field{{Key: "PORT", Type: schema.TypeInt, Required: true}}}
	path := writeSchema(t, s)
	loaded, err := schema.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded.Fields) != 1 || loaded.Fields[0].Key != "PORT" {
		t.Errorf("unexpected fields: %+v", loaded.Fields)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := schema.Load("/nonexistent/schema.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCheck_NoViolations(t *testing.T) {
	s := &schema.Schema{Fields: []schema.Field{
		{Key: "PORT", Type: schema.TypeInt, Required: true},
		{Key: "DEBUG", Type: schema.TypeBool, Required: false},
	}}
	env := map[string]string{"PORT": "8080", "DEBUG": "true"}
	if v := s.Check(env); len(v) != 0 {
		t.Errorf("expected no violations, got %v", v)
	}
}

func TestCheck_RequiredMissing(t *testing.T) {
	s := &schema.Schema{Fields: []schema.Field{
		{Key: "SECRET", Type: schema.TypeString, Required: true},
	}}
	v := s.Check(map[string]string{})
	if len(v) != 1 || v[0].Key != "SECRET" {
		t.Errorf("expected violation for SECRET, got %v", v)
	}
}

func TestCheck_TypeMismatch(t *testing.T) {
	s := &schema.Schema{Fields: []schema.Field{
		{Key: "PORT", Type: schema.TypeInt, Required: true},
	}}
	v := s.Check(map[string]string{"PORT": "not-a-number"})
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "PORT" {
		t.Errorf("expected violation key PORT, got %s", v[0].Key)
	}
}

func TestCheck_PatternViolation(t *testing.T) {
	s := &schema.Schema{Fields: []schema.Field{
		{Key: "ENV", Type: schema.TypeString, Required: true, Pattern: "(production|staging|development)"},
	}}
	v := s.Check(map[string]string{"ENV": "local"})
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestCheck_URLType(t *testing.T) {
	s := &schema.Schema{Fields: []schema.Field{
		{Key: "API_URL", Type: schema.TypeURL, Required: true},
	}}
	if v := s.Check(map[string]string{"API_URL": "https://api.example.com"}); len(v) != 0 {
		t.Errorf("expected no violations, got %v", v)
	}
	if v := s.Check(map[string]string{"API_URL": "ftp://bad"}); len(v) == 0 {
		t.Error("expected violation for non-http URL")
	}
}
