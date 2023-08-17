package command

import (
	"errors"

	"github.com/urfave/cli"
	"oldnicke/etcd/client"
)

// NewRemoveCommand returns the CLI command for "rm".
func NewRemoveCommand() cli.Command {
	return cli.Command{
		Name:      "rm",
		Usage:     "remove a key or a directory",
		ArgsUsage: "<key>",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "dir", Usage: "removes the key if it is an empty directory or a key-value pair"},
			cli.BoolFlag{Name: "recursive, r", Usage: "removes the key and all child keys(if it is a directory)"},
			cli.StringFlag{Name: "with-value", Value: "", Usage: "previous value"},
			cli.IntFlag{Name: "with-index", Value: 0, Usage: "previous index"},
		},
		Action: func(c *cli.Context) error {
			rmCommandFunc(c, mustNewKeyAPI(c))
			return nil
		},
	}
}

// rmCommandFunc executes the "rm" command.
func rmCommandFunc(c *cli.Context, ki client.KeysAPI) {
	if len(c.Args()) == 0 {
		handleError(c, ExitBadArgs, errors.New("key required"))
	}
	key := c.Args()[0]
	recursive := c.Bool("recursive")
	dir := c.Bool("dir")
	prevValue := c.String("with-value")
	prevIndex := c.Int("with-index")

	ctx, cancel := contextWithTotalTimeout(c)
	resp, err := ki.Delete(ctx, key, &client.DeleteOptions{PrevIndex: uint64(prevIndex), PrevValue: prevValue, Dir: dir, Recursive: recursive})
	cancel()
	if err != nil {
		handleError(c, ExitServerError, err)
	}
	if !resp.Node.Dir || c.GlobalString("output") != "simple" {
		printResponseKey(resp, c.GlobalString("output"))
	}
}
