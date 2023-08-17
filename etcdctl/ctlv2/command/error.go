package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/oldnicke/etcd/client"
	"github.com/urfave/cli"
)

const (
	ExitSuccess = iota
	ExitBadArgs
	ExitBadConnection
	ExitBadAuth
	ExitServerError
	ExitClusterNotHealthy
)

func handleError(c *cli.Context, code int, err error) {
	if c.GlobalString("output") == "json" {
		if err, ok := err.(*client.Error); ok {
			b, err := json.Marshal(err)
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(os.Stderr, string(b))
			os.Exit(code)
		}
	}

	fmt.Fprintln(os.Stderr, "Error: ", err)
	if cerr, ok := err.(*client.ClusterError); ok {
		fmt.Fprintln(os.Stderr, cerr.Detail())
	}
	os.Exit(code)
}
