// Command node-modules-scan discovers node_modules under ~ and emits inventory JSONL.
package main

import (
	"os"

	"disk-usage-analyser/nmscan"
)

func main() {
	code, err := nmscan.RunCLI(os.Args[1:], nmscan.CLIOptions{})
	if err != nil {
		os.Exit(1)
	}
	os.Exit(code)
}