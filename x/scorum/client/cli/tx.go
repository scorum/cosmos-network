package cli

import (
	"fmt"

	addresscodec "cosmossdk.io/core/address"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(ac addresscodec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetSupervisorTxCmd(ac),
	)

	return cmd
}

func GetSupervisorTxCmd(ac addresscodec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "admin",
		Short:                      "Supervisor transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdBurn(),
		CmdMintGas(),
		CmdSubmitMintProposal(ac),
	)

	return cmd
}
