package simulation

import (
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyGasLimit),
			func(r *rand.Rand) string {
				return fmt.Sprintf(
					`{"int": "%d"}`,
					simtypes.RandIntBetween(r, 1, 500000),
				)
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyGasUnconditionedAmount),
			func(r *rand.Rand) string {
				return fmt.Sprintf(
					`{"int": "%d"}`,
					simtypes.RandIntBetween(r, 1, 500000),
				)
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyGasAdjustCoefficient),
			func(r *rand.Rand) string {
				return fmt.Sprintf(
					`{"dec": "%s"}`,
					sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 2000)), 3),
				)
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeySPWithdrawalTotalPeriods),
			func(r *rand.Rand) string {
				return fmt.Sprintf(
					`%d`,
					genSpWithdrawalTotalPeriods(r),
				)
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeySPWithdrawalPeriodDurationSeconds),
			func(r *rand.Rand) string {
				return fmt.Sprintf(
					`%d`,
					genSpWithdrawalPeriodDurationSeconds(r),
				)
			},
		),
	}
}
