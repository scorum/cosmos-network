package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/spf13/cobra"
)

func CmdRegisterAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-account [addr]",
		Short: "Broadcast message MsgRegisterAccount",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRegisterAccount(
				clientCtx.GetFromAddress().String(),
				args[0],
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
