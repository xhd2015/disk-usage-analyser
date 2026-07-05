// Command node-modules-cache-shared reads a node_modules inventory JSON and emits
// JSONL with pnpm_cache_shared and bun_cache_shared per entry as each completes.
package main

import (
	"os"

	"disk-usage-analyser/nmcacheshared"
)

func main() {
	code, err := nmcacheshared.RunCLI(os.Args[1:], nmcacheshared.CLIOptions{})
	if err != nil {
		os.Exit(1)
	}
	os.Exit(code)
}