package command

import (
	"github.com/oldnicke/etcd/client"
	"github.com/urfave/cli"
)

// NewSetDirCommand returns the CLI command for "setDir".
func NewSetDirCommand() cli.Command {
	return cli.Command{
		Name:      "setdir",
		Usage:     "create a new directory or update an existing directory TTL",
		ArgsUsage: "<key>",
		Flags: []cli.Flag{
			cli.IntFlag{Name: "ttl", Value: 0, Usage: "key time-to-live in seconds"},
		},
		Action: func(c *cli.Context) error {
			mkdirCommandFunc(c, mustNewKeyAPI(c), client.PrevIgnore)
			return nil
		},
	}
}
