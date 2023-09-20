package cli

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/spf13/cobra"
)

func CmdWithdrawSP() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-sp [recipient] [amount]",
		Short: "Broadcast message MsgWithdrawSP",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			recipient, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid recipient address")
			}

			amount, ok := sdkmath.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid amount")
			}

			msg := types.NewMsgWithdrawSP(
				clientCtx.GetFromAddress().String(),
				recipient.String(),
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
