package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

func TestGetParams(t *testing.T) {
	set, ctx := setupKeeper(t)
	params := types.DefaultParams()

	set.keeper.SetParams(ctx.Context, params)

	require.EqualValues(t, params, set.keeper.GetParams(ctx.Context))
}
