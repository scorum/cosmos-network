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

func TestMsgServer_Burn(t *testing.T) {
	set, ctx := setupKeeper(t)

	coin := sdk.NewCoin(types.SCRDenom, math.NewInt(1000))

	s := keeper.NewMsgServer(set.keeper)

	require.NoError(t, set.keeper.Mint(ctx.Context, set.supervisor, coin))

	require.True(t, coin.Equal(set.bankKeeper.GetBalance(ctx.Context, set.supervisor, types.SCRDenom)))

	_, err := s.Burn(ctx, &types.MsgBurn{
		Supervisor: set.supervisor.String(),
		Amount:     coin,
	})
	require.NoError(t, err)

	require.True(t, set.bankKeeper.GetBalance(ctx.Context, set.supervisor, types.SCRDenom).IsZero())
}

func TestMsgServer_Burn_NotSupervisor(t *testing.T) {
	set, ctx := setupKeeper(t)

	addr := sample.AccAddress()
	coin := sdk.NewCoin(types.SCRDenom, math.NewInt(1000))

	s := keeper.NewMsgServer(set.keeper)

	require.NoError(t, set.keeper.Mint(ctx.Context, addr, coin))

	require.True(t, coin.Equal(set.bankKeeper.GetBalance(ctx.Context, addr, types.SCRDenom)))

	_, err := s.Burn(ctx, &types.MsgBurn{
		Supervisor: addr.String(),
		Amount:     coin,
	})
	require.Error(t, err)
}
