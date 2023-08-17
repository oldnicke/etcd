package main

import (
	"fmt"
	"os"

	"oldnicke/etcd/tools/benchmark/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
