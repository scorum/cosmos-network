package keeper_test

import (
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testutil "github.com/scorum/cosmos-network/testutil/keeper"
	"github.com/scorum/cosmos-network/testutil/sample"
	"github.com/scorum/cosmos-network/x/aviatrix/keeper"
	nfttypes "github.com/scorum/cosmos-network/x/aviatrix/types"
	scorumkeeper "github.com/scorum/cosmos-network/x/scorum/keeper"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
)

type keeperSet struct {
	supervisor sdk.AccAddress

	keeper       keeper.Keeper
	nftKeeper    nfttypes.NftKeeper
	scorumKeeper scorumkeeper.Keeper
}

func setupKeeper(t testing.TB) (keeperSet, testutil.TestContext) {
	ctx := testutil.GetContext(t)

	set := keeperSet{
		supervisor: sample.AccAddress(),

		keeper:       testutil.AviatrixKeeper(t, ctx),
		nftKeeper:    testutil.NftKeeper(t, ctx),
		scorumKeeper: testutil.ScorumKeeper(t, ctx),
	}

	set.scorumKeeper.SetParams(ctx.Context, scorumtypes.Params{
		Supervisors:            []string{set.supervisor.String()},
		GasLimit:               math.NewInt(1000),
		GasUnconditionedAmount: math.NewInt(500),
		GasAdjustCoefficient:   math.LegacyNewDec(1),
	})

	return set, ctx
}
