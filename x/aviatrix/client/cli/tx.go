package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "aviatrix",
		Short:                      "aviatrix transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
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
		CmdCreatePlane(),
		CmdUpdatePlaneExperience(),
		CmdAdjustPlaneExperience(),
	)

	return cmd
}
