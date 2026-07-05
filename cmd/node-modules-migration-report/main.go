// Command node-modules-migration-report runs before/after shared analyse with migrate in between.
package main

import (
	"os"

	"disk-usage-analyser/nmpipeline"
)

func main() {
	code, err := nmpipeline.RunCLI(os.Args[1:], nmpipeline.CLIOptions{})
	if err != nil {
		os.Exit(1)
	}
	os.Exit(code)
}