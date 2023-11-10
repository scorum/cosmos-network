package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

func CmdCreatePlane() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-plane [id] [owner] [experience]",
		Short: "Broadcast message CreatePlane",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argId := args[0]
			argOwner := args[1]
			var argExperience uint64

			if len(args) == 3 {
				argExperience, err = strconv.ParseUint(args[2], 10, 64)
				if err != nil {
					return fmt.Errorf("invalid experience: %w", err)
				}
			}

			msg := types.NewMsgCreatePlane(
				clientCtx.GetFromAddress().String(),
				argId,
				argOwner,
				argExperience,
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
