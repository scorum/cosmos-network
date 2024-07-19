package simulation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

// GenerateGenesisState creates a randomized GenState of the module
func GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	scorumGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}

	if len(simState.Accounts) > 0 {
		scorumGenesis.Params.Supervisors = []string{simState.Accounts[0].Address.String()}
	}
	if len(accs) > 0 {
		scorumGenesis.Params.ValidatorsReward.PoolAddress = accs[0]
		scorumGenesis.Params.ValidatorsReward.BlockReward = sdk.Coin{Amount: sdk.NewInt(1), Denom: types.SCRDenom}
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&scorumGenesis)
}
