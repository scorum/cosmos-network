package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

func TestParamsQuery(t *testing.T) {
	set, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx.Context)

	params := types.DefaultParams()

	set.keeper.SetParams(ctx.Context, params)

	act, err := set.keeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, act)
}
