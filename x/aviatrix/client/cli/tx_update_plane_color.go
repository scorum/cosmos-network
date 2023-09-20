package cli

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"

	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

func CmdUpdatePlaneColor() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-plane-color [id] [color]",
		Short: "Broadcast message MsgUpdatePlaneColor",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argId := args[0]
			argColor := args[1]

			if err != nil {
				return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid experience (%s)", err.Error())
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdatePlaneColor(
				clientCtx.GetFromAddress().String(),
				argId,
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
