package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

func CmdCreatePlane() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-plane [id] [owner] [name] [color]",
		Short: "Broadcast message CreatePlane",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argId := args[0]
			argOwner := args[1]
			argName := args[2]
			argColor := args[3]

			msg := types.NewMsgCreatePlane(
				clientCtx.GetFromAddress().String(),
				argId,
				argOwner,
				argName,
				argColor,
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
