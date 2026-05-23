package graph

import (
	"strings"
	"testing"
)

func TestBuild_NoRefs(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	g := Build(env)
	if len(g.Nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(g.Nodes))
	}
	if len(g.Nodes["HOST"].Refs) != 0 {
		t.Errorf("expected no refs for HOST")
	}
}

func TestBuild_ExtractsRefs(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "http://${HOST}:${PORT}",
		"HOST":     "localhost",
		"PORT":     "8080",
	}
	g := Build(env)
	refs := g.Nodes["BASE_URL"].Refs
	if len(refs) != 2 {
		t.Fatalf("expected 2 refs, got %d: %v", len(refs), refs)
	}
}

func TestMissingRefs_NoneWhenAllDefined(t *testing.T) {
	env := map[string]string{
		"A": "${B}",
		"B": "value",
	}
	g := Build(env)
	missing := g.MissingRefs()
	if len(missing) != 0 {
		t.Errorf("expected no missing refs, got %v", missing)
	}
}

func TestMissingRefs_DetectsUndefined(t *testing.T) {
	env := map[string]string{
		"DSN": "postgres://${DB_USER}:${DB_PASS}@${DB_HOST}/db",
	}
	g := Build(env)
	missing := g.MissingRefs()
	if len(missing) != 3 {
		t.Fatalf("expected 3 missing refs, got %v", missing)
	}
	if missing[0] != "DB_HOST" || missing[1] != "DB_PASS" || missing[2] != "DB_USER" {
		t.Errorf("unexpected missing refs order: %v", missing)
	}
}

func TestCyclicKeys_NoCycle(t *testing.T) {
	env := map[string]string{
		"A": "${B}",
		"B": "plain",
	}
	g := Build(env)
	if cyc := g.CyclicKeys(); len(cyc) != 0 {
		t.Errorf("expected no cycles, got %v", cyc)
	}
}

func TestCyclicKeys_DetectsCycle(t *testing.T) {
	env := map[string]string{
		"A": "${B}",
		"B": "${A}",
	}
	g := Build(env)
	cyc := g.CyclicKeys()
	if len(cyc) != 2 {
		t.Fatalf("expected 2 cyclic keys, got %v", cyc)
	}
}

func TestDot_ContainsEdge(t *testing.T) {
	env := map[string]string{
		"URL": "http://${HOST}",
		"HOST": "localhost",
	}
	g := Build(env)
	dot := g.Dot()
	if !strings.Contains(dot, "digraph envdiff") {
		t.Error("expected DOT header")
	}
	if !strings.Contains(dot, `"URL" -> "HOST"`) {
		t.Errorf("expected edge URL->HOST in DOT output:\n%s", dot)
	}
}

func TestDot_EmptyGraph(t *testing.T) {
	g := Build(map[string]string{})
	dot := g.Dot()
	if !strings.HasPrefix(dot, "digraph") {
		t.Error("expected valid DOT for empty graph")
	}
}
