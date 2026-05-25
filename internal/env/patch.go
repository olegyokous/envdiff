package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// PatchOp represents a single patch operation.
type PatchOp struct {
	Key    string
	Value  string
	Action string // "set", "delete", "rename"
	NewKey string // used for rename
}

// PatchResult holds the outcome of a patch operation.
type PatchResult struct {
	Op      PatchOp
	Applied bool
	Reason  string
}

// PatchOptions configures the Patch operation.
type PatchOptions struct {
	Source  string
	Output  string
	DryRun  bool
	Ops     []PatchOp
}

// Patch applies a sequence of operations to a .env file.
func Patch(opts PatchOptions) ([]PatchResult, error) {
	env, err := parser.ParseFile(opts.Source)
	if err != nil {
		return nil, fmt.Errorf("patch: parse %s: %w", opts.Source, err)
	}

	var results []PatchResult
	for _, op := range opts.Ops {
		result := applyOp(env, op)
		results = append(results, result)
	}

	if opts.DryRun {
		return results, nil
	}

	dest := opts.Output
	if dest == "" {
		dest = opts.Source
	}
	if err := writePatchFile(dest, env); err != nil {
		return nil, fmt.Errorf("patch: write %s: %w", dest, err)
	}
	return results, nil
}

func applyOp(env map[string]string, op PatchOp) PatchResult {
	switch op.Action {
	case "set":
		env[op.Key] = op.Value
		return PatchResult{Op: op, Applied: true}
	case "delete":
		if _, ok := env[op.Key]; !ok {
			return PatchResult{Op: op, Applied: false, Reason: "key not found"}
		}
		delete(env, op.Key)
		return PatchResult{Op: op, Applied: true}
	case "rename":
		val, ok := env[op.Key]
		if !ok {
			return PatchResult{Op: op, Applied: false, Reason: "key not found"}
		}
		if _, exists := env[op.NewKey]; exists {
			return PatchResult{Op: op, Applied: false, Reason: "target key already exists"}
		}
		env[op.NewKey] = val
		delete(env, op.Key)
		return PatchResult{Op: op, Applied: true}
	default:
		return PatchResult{Op: op, Applied: false, Reason: "unknown action"}
	}
}

func writePatchFile(path string, env map[string]string) error {
	var sb strings.Builder
	for k, v := range env {
		if strings.ContainsAny(v, " \t") {
			v = `"` + v + `"`
		}
		sb.WriteString(k + "=" + v + "\n")
	}
	return os.WriteFile(path, []byte(sb.String()), 0644)
}
