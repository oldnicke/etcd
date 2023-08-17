package command

import (
	"fmt"

	"go.etcd.io/etcd/version"

	"github.com/spf13/cobra"
)

// NewVersionCommand prints out the version of etcd.
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the version of etcdctl",
		Run:   versionCommandFunc,
	}
}

func versionCommandFunc(cmd *cobra.Command, args []string) {
	fmt.Println("etcdctl version:", version.Version)
	fmt.Println("API version:", version.APIVersion)
}
