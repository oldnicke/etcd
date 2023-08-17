package command

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"oldnicke/etcd/mvcc/backend"
)

var (
	defragDataDir string
)

// NewDefragCommand returns the cobra command for "Defrag".
func NewDefragCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "defrag",
		Short: "Defragments the storage of the etcd members with given endpoints",
		Run:   defragCommandFunc,
	}
	cmd.PersistentFlags().BoolVar(&epClusterEndpoints, "cluster", false, "use all endpoints from the cluster member list")
	cmd.Flags().StringVar(&defragDataDir, "data-dir", "", "Optional. If present, defragments a data directory not in use by etcd.")
	return cmd
}

func defragCommandFunc(cmd *cobra.Command, args []string) {
	if len(defragDataDir) > 0 {
		err := defragData(defragDataDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to defragment etcd data[%s] (%v)\n", defragDataDir, err)
			os.Exit(ExitError)
		}
		return
	}

	failures := 0
	c := mustClientFromCmd(cmd)
	for _, ep := range endpointsFromCluster(cmd) {
		ctx, cancel := commandCtx(cmd)
		_, err := c.Defragment(ctx, ep)
		cancel()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to defragment etcd member[%s] (%v)\n", ep, err)
			failures++
		} else {
			fmt.Printf("Finished defragmenting etcd member[%s]\n", ep)
		}
	}

	if failures != 0 {
		os.Exit(ExitError)
	}
}

func defragData(dataDir string) error {
	var be backend.Backend

	bch := make(chan struct{})
	dbDir := filepath.Join(dataDir, "member", "snap", "db")
	go func() {
		defer close(bch)
		be = backend.NewDefaultBackend(dbDir)

	}()
	select {
	case <-bch:
	case <-time.After(time.Second):
		fmt.Fprintf(os.Stderr, "waiting for etcd to close and release its lock on %q. "+
			"To defrag a running etcd instance, omit --data-dir.\n", dbDir)
		<-bch
	}
	return be.Defrag()
}
