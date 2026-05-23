package schema

import (
	"encoding/json"
	"fmt"
	"io"
)

// ViolationReport groups violations by environment name.
type ViolationReport struct {
	Env        string      `json:"env"`
	Violations []Violation `json:"violations"`
	OK         bool        `json:"ok"`
}

// WriteText writes a human-readable report of violations to w.
func WriteText(w io.Writer, reports []ViolationReport) {
	for _, r := range reports {
		if r.OK {
			fmt.Fprintf(w, "[%s] schema OK\n", r.Env)
			continue
		}
		fmt.Fprintf(w, "[%s] %d violation(s):\n", r.EnV, len(r.Violations))
		for _, v := range r.Violations {
			fmt.Fprintf(w, "  %-30s %s\n", v.Key, v.Message)
		}
	}
}

// WriteJSON writes a JSON array of ViolationReports to w.
func WriteJSON(w io.Writer, reports []ViolationReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(reports)
}

// Build constructs ViolationReports for each named env map.
func Build(s *Schema, envs map[string]map[string]string) []ViolationReport {
	reports := make([]ViolationReport, 0, len(envs))
	for name, env := range envs {
		violations := s.Check(env)
		reports = append(reports, ViolationReport{
			Env:        name,
			Violations: violations,
			OK:         len(violations) == 0,
		})
	}
	return reports
}

// HasViolations returns true if any report contains violations.
func HasViolations(reports []ViolationReport) bool {
	for _, r := range reports {
		if !r.OK {
			return true
		}
	}
	return false
}
