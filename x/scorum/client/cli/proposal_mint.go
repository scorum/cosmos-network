package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/address"

	addresscodec "cosmossdk.io/core/address"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/spf13/cobra"
)

const (
	FlagAuthority = "authority"
)

// CmdSubmitMintProposal implements a command handler for submitting a mint proposal transaction.
func CmdSubmitMintProposal(ac addresscodec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [receiver] [coin] [flags]",
		Args:  cobra.ExactArgs(2),
		Short: "Submit a mint proposal",
		Long:  "Submit a mint proposal along with an initial deposit.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal, err := cli.ReadGovPropFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			recipient := args[0]
			coin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return fmt.Errorf("invalid coin: %w", err)
			}

			authority, _ := cmd.Flags().GetString(FlagAuthority)
			if authority != "" {
				if _, err = ac.StringToBytes(authority); err != nil {
					return fmt.Errorf("invalid authority address: %w", err)
				}
			} else {
				authority = sdk.AccAddress(address.Module("gov")).String()
			}

			if err := proposal.SetMsgs([]sdk.Msg{
				&types.MsgMint{
					Authority: authority,
					Recipient: recipient,
					Amount:    coin,
				},
			}); err != nil {
				return fmt.Errorf("failed to create submit upgrade proposal message: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), proposal)
		},
	}

	cmd.Flags().String(FlagAuthority, "", "The address of the upgrade module authority (defaults to gov)")

	// add common proposal flags
	flags.AddTxFlagsToCmd(cmd)
	cli.AddGovPropFlagsToCmd(cmd)
	cmd.MarkFlagRequired(cli.FlagTitle)

	return cmd
}
