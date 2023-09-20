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

func CmdAdjustPlaneExperience() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "adjust-plane-experience [id] [amount]",
		Short: "Broadcast message MsgAdjustPlaneExperience",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argId := args[0]
			argAmount, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid amount (%s)", err.Error())
			}

			direction, err := cmd.Flags().GetString("direction")
			if err != nil {
				return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "failed to get direction flag")
			}

			switch direction {
			case "up":
			case "down":
				argAmount *= -1
			default:
				return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid direction")
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgAdjustPlaneExperience(
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
	cmd.Flags().String("direction", "up", "Increase or decrease plane experience. Possible values: up, down")

	return cmd
}
