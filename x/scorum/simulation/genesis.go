package simulation

import (
	"math/rand"

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
	scorumGenesis.Params.SpWithdrawalPeriodDurationSeconds = genSpWithdrawalPeriodDurationSeconds(simState.Rand)
	scorumGenesis.Params.SpWithdrawalTotalPeriods = genSpWithdrawalTotalPeriods(simState.Rand)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&scorumGenesis)
}

func genSpWithdrawalPeriodDurationSeconds(r *rand.Rand) uint32 {
	return uint32(r.Intn(5) + 1)
}

func genSpWithdrawalTotalPeriods(r *rand.Rand) uint32 {
	return uint32(r.Intn(100) + 1)
}
