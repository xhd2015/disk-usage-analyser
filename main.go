package main

import (
	"embed"
	"fmt"
	"os"

	"disk-usage-analyser/run"
	"disk-usage-analyser/server"
)

//go:embed disk-usage-analyser-react/dist
var distFS embed.FS

//go:embed disk-usage-analyser-react/template.html
var templateHTML string

func main() {
	server.Init(distFS, templateHTML)

	err := run.Run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
