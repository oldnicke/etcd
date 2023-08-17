package command

import (
	"fmt"

	"github.com/oldnicke/etcd/client"
	"github.com/urfave/cli"
)

func NewLsCommand() cli.Command {
	return cli.Command{
		Name:      "ls",
		Usage:     "retrieve a directory",
		ArgsUsage: "[key]",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "sort", Usage: "returns result in sorted order"},
			cli.BoolFlag{Name: "recursive, r", Usage: "returns all key names recursively for the given path"},
			cli.BoolFlag{Name: "p", Usage: "append slash (/) to directories"},
			cli.BoolFlag{Name: "quorum, q", Usage: "require quorum for get request"},
		},
		Action: func(c *cli.Context) error {
			lsCommandFunc(c, mustNewKeyAPI(c))
			return nil
		},
	}
}

// lsCommandFunc executes the "ls" command.
func lsCommandFunc(c *cli.Context, ki client.KeysAPI) {
	key := "/"
	if len(c.Args()) != 0 {
		key = c.Args()[0]
	}

	sort := c.Bool("sort")
	recursive := c.Bool("recursive")
	quorum := c.Bool("quorum")

	ctx, cancel := contextWithTotalTimeout(c)
	resp, err := ki.Get(ctx, key, &client.GetOptions{Sort: sort, Recursive: recursive, Quorum: quorum})
	cancel()
	if err != nil {
		handleError(c, ExitServerError, err)
	}

	printLs(c, resp)
}

// printLs writes a response out in a manner similar to the `ls` command in unix.
// Non-empty directories list their contents and files list their name.
func printLs(c *cli.Context, resp *client.Response) {
	if c.GlobalString("output") == "simple" {
		if !resp.Node.Dir {
			fmt.Println(resp.Node.Key)
		}
		for _, node := range resp.Node.Nodes {
			rPrint(c, node)
		}
	} else {
		// user wants JSON or extended output
		printResponseKey(resp, c.GlobalString("output"))
	}
}

// rPrint recursively prints out the nodes in the node structure.
func rPrint(c *cli.Context, n *client.Node) {
	if n.Dir && c.Bool("p") {
		fmt.Println(fmt.Sprintf("%v/", n.Key))
	} else {
		fmt.Println(n.Key)
	}

	for _, node := range n.Nodes {
		rPrint(c, node)
	}
}
