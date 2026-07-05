package main

import (
	"fmt"
	"os"

	"github.com/xhd2015/dot-pkgs/go-pkgs/npm"
	"github.com/xhd2015/xgo/support/cmd"
)

const reactDir = "disk-usage-analyser-react"

func main() {
	err := Handle(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func Handle(args []string) error {
	manager, err := npm.Resolve(reactDir, "auto")
	if err != nil {
		return err
	}

	if _, err := os.Stat(reactDir + "/node_modules"); err != nil {
		name, installArgs := npm.InstallCommand(manager, npm.InstallOptions{})
		if err := cmd.Debug().Dir(reactDir).Run(name, installArgs...); err != nil {
			return err
		}
	}

	return cmd.Debug().Dir(reactDir).Run(string(manager), "run", "build")
}