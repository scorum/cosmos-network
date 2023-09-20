package cli

import (
	"strconv"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"

	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

func CmdUpdatePlaneExperience() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-plane-experience [id] [amount]",
		Short: "Broadcast message MsgUpdatePlaneExperience",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argId := args[0]
			argAmount, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid amount (%s)", err.Error())
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdatePlaneExperience(
				clientCtx.GetFromAddress().String(),
				argId,
				argAmount,
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
