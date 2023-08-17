package command

import (
	"errors"
	"time"

	"github.com/urfave/cli"
	"oldnicke/etcd/client"
)

// NewUpdateDirCommand returns the CLI command for "updatedir".
func NewUpdateDirCommand() cli.Command {
	return cli.Command{
		Name:      "updatedir",
		Usage:     "update an existing directory",
		ArgsUsage: "<key> <value>",
		Flags: []cli.Flag{
			cli.IntFlag{Name: "ttl", Value: 0, Usage: "key time-to-live in seconds"},
		},
		Action: func(c *cli.Context) error {
			updatedirCommandFunc(c, mustNewKeyAPI(c))
			return nil
		},
	}
}

// updatedirCommandFunc executes the "updatedir" command.
func updatedirCommandFunc(c *cli.Context, ki client.KeysAPI) {
	if len(c.Args()) == 0 {
		handleError(c, ExitBadArgs, errors.New("key required"))
	}
	key := c.Args()[0]
	ttl := c.Int("ttl")
	ctx, cancel := contextWithTotalTimeout(c)
	resp, err := ki.Set(ctx, key, "", &client.SetOptions{TTL: time.Duration(ttl) * time.Second, Dir: true, PrevExist: client.PrevExist})
	cancel()
	if err != nil {
		handleError(c, ExitServerError, err)
	}
	if c.GlobalString("output") != "simple" {
		printResponseKey(resp, c.GlobalString("output"))
	}
}
