package command

import (
	"errors"
	"time"

	"github.com/urfave/cli"
	"go.etcd.io/etcd/client"
)

// NewMakeDirCommand returns the CLI command for "mkdir".
func NewMakeDirCommand() cli.Command {
	return cli.Command{
		Name:      "mkdir",
		Usage:     "make a new directory",
		ArgsUsage: "<key>",
		Flags: []cli.Flag{
			cli.IntFlag{Name: "ttl", Value: 0, Usage: "key time-to-live in seconds"},
		},
		Action: func(c *cli.Context) error {
			mkdirCommandFunc(c, mustNewKeyAPI(c), client.PrevNoExist)
			return nil
		},
	}
}

// mkdirCommandFunc executes the "mkdir" command.
func mkdirCommandFunc(c *cli.Context, ki client.KeysAPI, prevExist client.PrevExistType) {
	if len(c.Args()) == 0 {
		handleError(c, ExitBadArgs, errors.New("key required"))
	}

	key := c.Args()[0]
	ttl := c.Int("ttl")

	ctx, cancel := contextWithTotalTimeout(c)
	resp, err := ki.Set(ctx, key, "", &client.SetOptions{TTL: time.Duration(ttl) * time.Second, Dir: true, PrevExist: prevExist})
	cancel()
	if err != nil {
		handleError(c, ExitServerError, err)
	}
	if c.GlobalString("output") != "simple" {
		printResponseKey(resp, c.GlobalString("output"))
	}
}
