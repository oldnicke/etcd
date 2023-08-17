package command

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/oldnicke/etcd/client"

	"github.com/urfave/cli"
)

// NewWatchCommand returns the CLI command for "watch".
func NewWatchCommand() cli.Command {
	return cli.Command{
		Name:      "watch",
		Usage:     "watch a key for changes",
		ArgsUsage: "<key>",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "forever, f", Usage: "forever watch a key until CTRL+C"},
			cli.IntFlag{Name: "after-index", Value: 0, Usage: "watch after the given index"},
			cli.BoolFlag{Name: "recursive, r", Usage: "returns all values for key and child keys"},
		},
		Action: func(c *cli.Context) error {
			watchCommandFunc(c, mustNewKeyAPI(c))
			return nil
		},
	}
}

// watchCommandFunc executes the "watch" command.
func watchCommandFunc(c *cli.Context, ki client.KeysAPI) {
	if len(c.Args()) == 0 {
		handleError(c, ExitBadArgs, errors.New("key required"))
	}
	key := c.Args()[0]
	recursive := c.Bool("recursive")
	forever := c.Bool("forever")
	index := c.Int("after-index")

	stop := false
	w := ki.Watcher(key, &client.WatcherOptions{AfterIndex: uint64(index), Recursive: recursive})

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)

	go func() {
		<-sigch
		os.Exit(0)
	}()

	for !stop {
		resp, err := w.Next(context.TODO())
		if err != nil {
			handleError(c, ExitServerError, err)
		}
		if resp.Node.Dir {
			continue
		}
		if recursive {
			fmt.Printf("[%s] %s\n", resp.Action, resp.Node.Key)
		}

		printResponseKey(resp, c.GlobalString("output"))

		if !forever {
			stop = true
		}
	}
}
