package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
	"go.etcd.io/etcd/client"
)

func NewAuthCommands() cli.Command {
	return cli.Command{
		Name:  "auth",
		Usage: "overall auth controls",
		Subcommands: []cli.Command{
			{
				Name:      "enable",
				Usage:     "enable auth access controls",
				ArgsUsage: " ",
				Action:    actionAuthEnable,
			},
			{
				Name:      "disable",
				Usage:     "disable auth access controls",
				ArgsUsage: " ",
				Action:    actionAuthDisable,
			},
		},
	}
}

func actionAuthEnable(c *cli.Context) error {
	authEnableDisable(c, true)
	return nil
}

func actionAuthDisable(c *cli.Context) error {
	authEnableDisable(c, false)
	return nil
}

func mustNewAuthAPI(c *cli.Context) client.AuthAPI {
	hc := mustNewClient(c)

	if c.GlobalBool("debug") {
		fmt.Fprintf(os.Stderr, "Cluster-Endpoints: %s\n", strings.Join(hc.Endpoints(), ", "))
	}

	return client.NewAuthAPI(hc)
}

func authEnableDisable(c *cli.Context, enable bool) {
	if len(c.Args()) != 0 {
		fmt.Fprintln(os.Stderr, "No arguments accepted")
		os.Exit(1)
	}
	s := mustNewAuthAPI(c)
	ctx, cancel := contextWithTotalTimeout(c)
	var err error
	if enable {
		err = s.Enable(ctx)
	} else {
		err = s.Disable(ctx)
	}
	cancel()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if enable {
		fmt.Println("Authentication Enabled")
	} else {
		fmt.Println("Authentication Disabled")
	}
}
