package command

import (
	"fmt"

	"github.com/oldnicke/etcd/clientv3"
	"github.com/spf13/cobra"
)

var (
	delPrefix  bool
	delPrevKV  bool
	delFromKey bool
)

// NewDelCommand returns the cobra command for "del".
func NewDelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "del [options] <key> [range_end]",
		Short: "Removes the specified key or range of keys [key, range_end)",
		Run:   delCommandFunc,
	}

	cmd.Flags().BoolVar(&delPrefix, "prefix", false, "delete keys with matching prefix")
	cmd.Flags().BoolVar(&delPrevKV, "prev-kv", false, "return deleted key-value pairs")
	cmd.Flags().BoolVar(&delFromKey, "from-key", false, "delete keys that are greater than or equal to the given key using byte compare")
	return cmd
}

// delCommandFunc executes the "del" command.
func delCommandFunc(cmd *cobra.Command, args []string) {
	key, opts := getDelOp(args)
	ctx, cancel := commandCtx(cmd)
	resp, err := mustClientFromCmd(cmd).Delete(ctx, key, opts...)
	cancel()
	if err != nil {
		ExitWithError(ExitError, err)
	}
	display.Del(*resp)
}

func getDelOp(args []string) (string, []clientv3.OpOption) {
	if len(args) == 0 || len(args) > 2 {
		ExitWithError(ExitBadArgs, fmt.Errorf("del command needs one argument as key and an optional argument as range_end"))
	}

	if delPrefix && delFromKey {
		ExitWithError(ExitBadArgs, fmt.Errorf("`--prefix` and `--from-key` cannot be set at the same time, choose one"))
	}

	opts := []clientv3.OpOption{}
	key := args[0]
	if len(args) > 1 {
		if delPrefix || delFromKey {
			ExitWithError(ExitBadArgs, fmt.Errorf("too many arguments, only accept one argument when `--prefix` or `--from-key` is set"))
		}
		opts = append(opts, clientv3.WithRange(args[1]))
	}

	if delPrefix {
		if len(key) == 0 {
			key = "\x00"
			opts = append(opts, clientv3.WithFromKey())
		} else {
			opts = append(opts, clientv3.WithPrefix())
		}
	}
	if delPrevKV {
		opts = append(opts, clientv3.WithPrevKV())
	}

	if delFromKey {
		if len(key) == 0 {
			key = "\x00"
		}
		opts = append(opts, clientv3.WithFromKey())
	}

	return key, opts
}
