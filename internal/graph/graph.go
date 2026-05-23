// Package graph builds a dependency graph between .env keys based on
// value references (e.g. KEY=${OTHER}) and detects cycles or missing refs.
package graph

import (
	"fmt"
	"regexp"
	"sort"
)

// refPattern matches ${VAR} or $VAR style references inside values.
var refPattern = regexp.MustCompile(`\$\{?([A-Z_][A-Z0-9_]*)\}?`)

// Node represents a single key and the keys it references.
type Node struct {
	Key  string
	Refs []string
}

// Graph holds the full dependency graph for one environment.
type Graph struct {
	Nodes map[string]*Node
}

// Build constructs a Graph from a parsed env map.
func Build(env map[string]string) *Graph {
	g := &Graph{Nodes: make(map[string]*Node, len(env))}
	for k, v := range env {
		refs := extractRefs(v)
		g.Nodes[k] = &Node{Key: k, Refs: refs}
	}
	return g
}

// MissingRefs returns keys that are referenced but not defined in the env.
func (g *Graph) MissingRefs() []string {
	seen := make(map[string]bool)
	var missing []string
	for _, node := range g.Nodes {
		for _, ref := range node.Refs {
			if _, defined := g.Nodes[ref]; !defined && !seen[ref] {
				missing = append(missing, ref)
				seen[ref] = true
			}
		}
	}
	sort.Strings(missing)
	return missing
}

// CyclicKeys returns keys that participate in a reference cycle.
func (g *Graph) CyclicKeys() []string {
	visited := make(map[string]bool)
	inStack := make(map[string]bool)
	cyclic := make(map[string]bool)

	var dfs func(key string)
	dfs = func(key string) {
		visited[key] = true
		inStack[key] = true
		node, ok := g.Nodes[key]
		if ok {
			for _, ref := range node.Refs {
				if !visited[ref] {
					dfs(ref)
				} else if inStack[ref] {
					cyclic[ref] = true
					cyclic[key] = true
				}
			}
		}
		inStack[key] = false
	}

	for k := range g.Nodes {
		if !visited[k] {
			dfs(k)
		}
	}

	keys := make([]string, 0, len(cyclic))
	for k := range cyclic {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Dot returns a Graphviz DOT representation of the graph.
func (g *Graph) Dot() string {
	out := "digraph envdiff {\n"
	keys := make([]string, 0, len(g.Nodes))
	for k := range g.Nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		for _, ref := range g.Nodes[k].Refs {
			out += fmt.Sprintf("  %q -> %q;\n", k, ref)
		}
	}
	out += "}\n"
	return out
}

func extractRefs(value string) []string {
	matches := refPattern.FindAllStringSubmatch(value, -1)
	seen := make(map[string]bool)
	var refs []string
	for _, m := range matches {
		if !seen[m[1]] {
			refs = append(refs, m[1])
			seen[m[1]] = true
		}
	}
	return refs
}
