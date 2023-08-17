package main

import (
	"fmt"
	"os"

	"go.etcd.io/etcd/tools/benchmark/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
