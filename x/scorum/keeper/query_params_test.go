package keeper_test

import (
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scorum/cosmos-network/testutil/sample"
	"github.com/scorum/cosmos-network/x/scorum/keeper"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/stretchr/testify/require"
)

func TestQuery_GetParams(t *testing.T) {
	set, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx.Context)

	params := types.NewParams(
		[]string{sample.AccAddress().String()},
		math.NewInt(1000),
		math.NewInt(500),
		math.LegacyNewDec(1),
	)
	params.ValidatorsReward = types.ValidatorsRewardParams{
		PoolAddress: "",
		BlockReward: sdk.Coin{
			Denom:  "",
			Amount: math.ZeroInt(),
		},
	}

	set.keeper.SetParams(ctx.Context, params)
	srv := keeper.NewQueryServer(set.keeper)

	response, err := srv.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
