package main

import (
	"fmt"
	"os"

	"github.com/xhd2015/xgo/support/cmd"
)

func main() {
	err := Handle(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func Handle(args []string) error {
	return cmd.Debug().Run("go", append([]string{"run", ".", "--dev"}, args...)...)
}
