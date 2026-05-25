package env

import (
	"fmt"
	"io"
	"strings"
)

// PatchCmdOptions holds CLI-level options for the patch command.
type PatchCmdOptions struct {
	Source  string
	Output  string
	DryRun  bool
	Format  string   // "text" or "json"
	RawOps  []string // "set:KEY=VALUE", "delete:KEY", "rename:OLD:NEW"
}

// RunPatch parses raw operation strings and applies them.
func RunPatch(opts PatchCmdOptions, w io.Writer) error {
	ops, err := parseRawOps(opts.RawOps)
	if err != nil {
		return fmt.Errorf("patch: invalid ops: %w", err)
	}

	results, err := Patch(PatchOptions{
		Source: opts.Source,
		Output: opts.Output,
		DryRun: opts.DryRun,
		Ops:    ops,
	})
	if err != nil {
		return err
	}

	switch strings.ToLower(opts.Format) {
	case "json":
		return WritePatchJSON(w, results)
	default:
		WritePatchText(w, results, opts.Source)
		return nil
	}
}

func parseRawOps(raw []string) ([]PatchOp, error) {
	var ops []PatchOp
	for _, s := range raw {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid op %q: expected action:args", s)
		}
		action := strings.ToLower(parts[0])
		args := parts[1]
		switch action {
		case "set":
			kv := strings.SplitN(args, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("set op requires KEY=VALUE, got %q", args)
			}
			ops = append(ops, PatchOp{Action: "set", Key: kv[0], Value: kv[1]})
		case "delete":
			ops = append(ops, PatchOp{Action: "delete", Key: args})
		case "rename":
			keys := strings.SplitN(args, ":", 2)
			if len(keys) != 2 {
				return nil, fmt.Errorf("rename op requires OLD:NEW, got %q", args)
			}
			ops = append(ops, PatchOp{Action: "rename", Key: keys[0], NewKey: keys[1]})
		default:
			return nil, fmt.Errorf("unknown action %q", action)
		}
	}
	return ops, nil
}
