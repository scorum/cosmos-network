package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/google/uuid"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/spf13/cobra"
)

func CmdStopSPWithdrawal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop-sp-withdrawal [id]",
		Short: "Broadcast message MsgStopSPWithdrawal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := uuid.Parse(args[0])
			if err != nil {
				return fmt.Errorf("invalid id: must be uuid")
			}

			msg := types.NewMsgStopSPWithdrawal(
				clientCtx.GetFromAddress().String(),
				id.String(),
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
