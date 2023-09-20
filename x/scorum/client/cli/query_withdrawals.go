package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/spf13/cobra"
)

func CmdQueryWithdrawals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdrawals [owner]",
		Short: "list account's withdrawals",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid owner address: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.ListWithdrawals(context.Background(), &types.QueryWithdrawalsRequest{
				Owner: owner.String(),
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
