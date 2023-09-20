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

func CmdUpdatePlaneName() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-plane-name [id] [name]",
		Short: "Broadcast message MsgUpdatePlaneName",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argId := args[0]
			argName := args[1]

			if err != nil {
				return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid experience (%s)", err.Error())
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdatePlaneName(
				clientCtx.GetFromAddress().String(),
				argId,
				argName,
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
