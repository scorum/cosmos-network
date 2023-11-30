package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/os"
)

func CmdSubmitMintProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a mint proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a mint proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.
Example:
$ %s tx gov submit-proposal mint <path/to/proposal.json> --from=<key_or_address>
Where proposal.json contains:
{
  "title": "Mint burned on old blockchain coins",
  "description": "Emission in the modern blockchain!",
  "recipient": "scorum16rlsek5yak6avnjpuatw6psxg246nzlwruaet5",
  "amount": 
  	{
			"denom": "nscr",
			"amount": "100000000000"
	},
  "deposit": "1000scr"
}
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			proposal, err := ParseMintProposalJSON(clientCtx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			content := types.NewMintProposal(
				proposal.Title, proposal.Description, proposal.Recipient, proposal.Amount,
			)

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			msg, err := govv1beta1.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

type MintProposalJSON struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Deposit     string `json:"deposit" yaml:"deposit"`

	Recipient sdk.AccAddress `json:"recipient" yaml:"recipient"`
	Amount    sdk.Coin       `json:"amount" yaml:"amount"`
}

func ParseMintProposalJSON(cdc *codec.LegacyAmino, proposalFile string) (MintProposalJSON, error) {
	proposal := MintProposalJSON{}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
