package main

import (
	"fmt"
	"os"

	"github.com/yourorg/envdiff/internal/parser"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: envdiff <base.env> <compare.env>")
		os.Exit(1)
	}

	basePath := os.Args[1]
	cmpPath := os.Args[2]

	baseEnv, err := parser.ParseFile(basePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing base file: %v\n", err)
		os.Exit(1)
	}

	cmpEnv, err := parser.ParseFile(cmpPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing compare file: %v\n", err)
		os.Exit(1)
	}

	missing := []string{}
	extra := []string{}

	for k := range baseEnv {
		if _, ok := cmpEnv[k]; !ok {
			missing = append(missing, k)
		}
	}

	for k := range cmpEnv {
		if _, ok := baseEnv[k]; !ok {
			extra = append(extra, k)
		}
	}

	if len(missing) == 0 && len(extra) == 0 {
		fmt.Println("No differences found.")
		return
	}

	for _, k := range missing {
		fmt.Printf("MISSING  %s\n", k)
	}
	for _, k := range extra {
		fmt.Printf("EXTRA    %s\n", k)
	}

	os.Exit(1)
}
