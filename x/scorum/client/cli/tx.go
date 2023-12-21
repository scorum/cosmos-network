package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdConvertSCR2SP(),
		CmdWithdrawSP(),
		CmdStopSPWithdrawal(),
		GetSupervisorTxCmd(),
	)

	return cmd
}

func GetSupervisorTxCmd() *cobra.Command {
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
	)

	return cmd
}
