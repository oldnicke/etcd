package command

import (
	"errors"

	"github.com/urfave/cli"
	"oldnicke/etcd/client"
)

// NewRemoveDirCommand returns the CLI command for "rmdir".
func NewRemoveDirCommand() cli.Command {
	return cli.Command{
		Name:      "rmdir",
		Usage:     "removes the key if it is an empty directory or a key-value pair",
		ArgsUsage: "<key>",
		Action: func(c *cli.Context) error {
			rmdirCommandFunc(c, mustNewKeyAPI(c))
			return nil
		},
	}
}

// rmdirCommandFunc executes the "rmdir" command.
func rmdirCommandFunc(c *cli.Context, ki client.KeysAPI) {
	if len(c.Args()) == 0 {
		handleError(c, ExitBadArgs, errors.New("key required"))
	}
	key := c.Args()[0]

	ctx, cancel := contextWithTotalTimeout(c)
	resp, err := ki.Delete(ctx, key, &client.DeleteOptions{Dir: true})
	cancel()
	if err != nil {
		handleError(c, ExitServerError, err)
	}

	if !resp.Node.Dir || c.GlobalString("output") != "simple" {
		printResponseKey(resp, c.GlobalString("output"))
	}
}
