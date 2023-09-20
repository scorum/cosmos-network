package cli

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/spf13/cobra"
)

func CmdConvertSCR2SP() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert-scr2sp [amount]",
		Short: "Broadcast message MsgConvertSCR2SP",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, ok := sdkmath.NewIntFromString(args[0])
			if !ok {
				return fmt.Errorf("invalid amount")
			}

			msg := types.NewMsgConvertSCR2SP(
				clientCtx.GetFromAddress().String(),
				amount,
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
