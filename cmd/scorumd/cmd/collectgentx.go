package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/scorum/cosmos-network/app"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/spf13/cobra"
)

func CollectGenTxsCmd(genBalIterator banktypes.GenesisBalancesIterator, defaultNodeHome string) *cobra.Command {
	gentxModule, ok := app.ModuleBasics[genutiltypes.ModuleName].(genutil.AppModuleBasic)
	if !ok {
		panic(fmt.Errorf("expected %s module to be an instance of type %T", genutiltypes.ModuleName, genutil.AppModuleBasic{}))
	}

	cmd := genutilcli.CollectGenTxsCmd(genBalIterator, defaultNodeHome, gentxModule.GenTxValidator)

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		clientCtx := client.GetClientContextFromCmd(cmd)
		cdc := clientCtx.Codec

		genFile := server.GetServerContextFromCmd(cmd).Config.GenesisFile()
		appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
		if err != nil {
			return fmt.Errorf("failed to unmarshal genesis state: %w", err)
		}

		bankGenState := banktypes.GetGenesisStateFromAppState(cdc, appState)
		bankGenState.SendEnabled = append(
			bankGenState.SendEnabled,
			banktypes.SendEnabled{Denom: scorumtypes.GasDenom, Enabled: false},
		)

		bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
		if err != nil {
			return fmt.Errorf("failed to marshal bank genesis state: %w", err)
		}
		appState[banktypes.ModuleName] = bankGenStateBz

		distrGenState := distrtypes.DefaultGenesisState()
		distrGenState.Params = distrtypes.Params{
			CommunityTax:        types.ZeroDec(),
			BaseProposerReward:  distrGenState.Params.BaseProposerReward,
			BonusProposerReward: distrGenState.Params.BonusProposerReward,
			WithdrawAddrEnabled: true,
		}
		distrGenStateBz, err := cdc.MarshalJSON(distrGenState)
		if err != nil {
			return fmt.Errorf("failed to marshal distr genesis state: %w", err)
		}
		appState[distrtypes.ModuleName] = distrGenStateBz

		mintGenState := minttypes.DefaultGenesisState()
		mintGenState.Params = minttypes.Params{
			MintDenom:           scorumtypes.SCRDenom,
			InflationRateChange: types.ZeroDec(),
			InflationMax:        types.ZeroDec(),
			InflationMin:        types.ZeroDec(),
			GoalBonded:          mintGenState.Params.GoalBonded,
			BlocksPerYear:       mintGenState.Params.BlocksPerYear,
		}
		mintGenStateBz, err := cdc.MarshalJSON(mintGenState)
		if err != nil {
			return fmt.Errorf("failed to marshal mint genesis state: %w", err)
		}
		appState[minttypes.ModuleName] = mintGenStateBz

		appStateJSON, err := json.Marshal(appState)
		if err != nil {
			return fmt.Errorf("failed to marshal application genesis state: %w", err)
		}

		genDoc.AppState = appStateJSON
		return genutil.ExportGenesisFile(genDoc, genFile)
	}

	return cmd
}
