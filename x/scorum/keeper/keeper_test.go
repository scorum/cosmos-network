package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	testutil "github.com/scorum/cosmos-network/testutil/keeper"
	"github.com/scorum/cosmos-network/testutil/sample"
	"github.com/scorum/cosmos-network/x/scorum/keeper"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

type keeperSet struct {
	supervisor sdk.AccAddress

	keeper        keeper.Keeper
	accountKeeper accountkeeper.AccountKeeper
	bankKeeper    bankkeeper.Keeper
}

func setupKeeper(t testing.TB) (keeperSet, testutil.TestContext) {
	ctx := testutil.GetContext(t)
	set := keeperSet{
		supervisor: sample.AccAddress(),

		keeper:        testutil.ScorumKeeper(t, ctx),
		accountKeeper: testutil.AccountKeeper(t, ctx),
		bankKeeper:    testutil.BankKeeper(t, ctx),
	}

	set.keeper.SetParams(ctx.Context, types.Params{
		Supervisors:            []string{set.supervisor.String()},
		GasLimit:               sdk.IntProto{Int: sdk.NewInt(1000)},
		GasUnconditionedAmount: sdk.IntProto{Int: sdk.NewInt(500)},
		GasAdjustCoefficient:   sdk.DecProto{Dec: sdk.NewDec(1)},
	})

	return set, ctx
}
