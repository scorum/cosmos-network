package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

// GenerateGenesisState creates a randomized GenState of the module
func GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	aviatrixGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&aviatrixGenesis)
}
