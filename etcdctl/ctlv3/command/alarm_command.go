package command

import (
	"fmt"

	"github.com/spf13/cobra"
	v3 "oldnicke/etcd/clientv3"
)

// NewAlarmCommand returns the cobra command for "alarm".
func NewAlarmCommand() *cobra.Command {
	ac := &cobra.Command{
		Use:   "alarm <subcommand>",
		Short: "Alarm related commands",
	}

	ac.AddCommand(NewAlarmDisarmCommand())
	ac.AddCommand(NewAlarmListCommand())

	return ac
}

func NewAlarmDisarmCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "disarm",
		Short: "Disarms all alarms",
		Run:   alarmDisarmCommandFunc,
	}
	return &cmd
}

// alarmDisarmCommandFunc executes the "alarm disarm" command.
func alarmDisarmCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		ExitWithError(ExitBadArgs, fmt.Errorf("alarm disarm command accepts no arguments"))
	}
	ctx, cancel := commandCtx(cmd)
	resp, err := mustClientFromCmd(cmd).AlarmDisarm(ctx, &v3.AlarmMember{})
	cancel()
	if err != nil {
		ExitWithError(ExitError, err)
	}
	display.Alarm(*resp)
}

func NewAlarmListCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list",
		Short: "Lists all alarms",
		Run:   alarmListCommandFunc,
	}
	return &cmd
}

// alarmListCommandFunc executes the "alarm list" command.
func alarmListCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		ExitWithError(ExitBadArgs, fmt.Errorf("alarm list command accepts no arguments"))
	}
	ctx, cancel := commandCtx(cmd)
	resp, err := mustClientFromCmd(cmd).AlarmList(ctx)
	cancel()
	if err != nil {
		ExitWithError(ExitError, err)
	}
	display.Alarm(*resp)
}
