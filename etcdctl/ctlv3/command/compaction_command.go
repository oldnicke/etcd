package command

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"oldnicke/etcd/clientv3"
)

var compactPhysical bool

// NewCompactionCommand returns the cobra command for "compaction".
func NewCompactionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compaction [options] <revision>",
		Short: "Compacts the event history in etcd",
		Run:   compactionCommandFunc,
	}
	cmd.Flags().BoolVar(&compactPhysical, "physical", false, "'true' to wait for compaction to physically remove all old revisions")
	return cmd
}

// compactionCommandFunc executes the "compaction" command.
func compactionCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("compaction command needs 1 argument"))
	}

	rev, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		ExitWithError(ExitError, err)
	}

	var opts []clientv3.CompactOption
	if compactPhysical {
		opts = append(opts, clientv3.WithCompactPhysical())
	}

	c := mustClientFromCmd(cmd)
	ctx, cancel := commandCtx(cmd)
	_, cerr := c.Compact(ctx, rev, opts...)
	cancel()
	if cerr != nil {
		ExitWithError(ExitError, cerr)
	}
	fmt.Println("compacted revision", rev)
}
