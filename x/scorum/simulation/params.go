package simulation

import (
	"fmt"
	"math/rand"

	"cosmossdk.io/math"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

func ParamChanges(r *rand.Rand) []simtypes.LegacyParamChange {
	return []simtypes.LegacyParamChange{
		simulation.NewSimLegacyParamChange(types.ModuleName, string(types.KeyGasLimit),
			func(r *rand.Rand) string {
				return fmt.Sprintf(
					`{"int": "%d"}`,
					simtypes.RandIntBetween(r, 1, 500000),
				)
			},
		),
		simulation.NewSimLegacyParamChange(types.ModuleName, string(types.KeyGasUnconditionedAmount),
			func(r *rand.Rand) string {
				return fmt.Sprintf(
					`{"int": "%d"}`,
					simtypes.RandIntBetween(r, 1, 500000),
				)
			},
		),
		simulation.NewSimLegacyParamChange(types.ModuleName, string(types.KeyGasAdjustCoefficient),
			func(r *rand.Rand) string {
				return fmt.Sprintf(
					`{"dec": "%s"}`,
					math.LegacyNewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 2000)), 3),
				)
			},
		),
	}
}
