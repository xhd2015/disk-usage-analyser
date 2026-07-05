// Command node-modules-migrate-pnpm removes node_modules and switches projects to pnpm via corepack.
package main

import (
	"os"

	"disk-usage-analyser/nmmigrate"
)

func main() {
	code, err := nmmigrate.RunCLI(os.Args[1:], nmmigrate.CLIOptions{})
	if err != nil {
		os.Exit(1)
	}
	os.Exit(code)
}