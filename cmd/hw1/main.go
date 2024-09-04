package main

import (
	"fmt"
	"homework/cmd/hw1/cmds"
	"os"
)

func main() {
	if err := cmds.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
