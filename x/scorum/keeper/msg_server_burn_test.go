package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scorum/cosmos-network/testutil/sample"
	"github.com/scorum/cosmos-network/x/scorum/keeper"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServer_Burn(t *testing.T) {
	set, ctx := setupKeeper(t)

	coin := sdk.NewCoin(types.SCRDenom, sdk.NewInt(1000))

	s := keeper.NewMsgServer(set.keeper)

	require.NoError(t, set.bankKeeper.MintCoins(ctx.Context, types.ModuleName, sdk.NewCoins(coin)))
	require.NoError(t, set.bankKeeper.SendCoinsFromModuleToAccount(ctx.Context, types.ModuleName, set.supervisor, sdk.NewCoins(coin)))

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
	coin := sdk.NewCoin(types.SCRDenom, sdk.NewInt(1000))

	s := keeper.NewMsgServer(set.keeper)

	require.NoError(t, set.bankKeeper.MintCoins(ctx.Context, types.ModuleName, sdk.NewCoins(coin)))
	require.NoError(t, set.bankKeeper.SendCoinsFromModuleToAccount(ctx.Context, types.ModuleName, addr, sdk.NewCoins(coin)))

	require.True(t, coin.Equal(set.bankKeeper.GetBalance(ctx.Context, addr, types.SCRDenom)))

	_, err := s.Burn(ctx, &types.MsgBurn{
		Supervisor: addr.String(),
		Amount:     coin,
	})
	require.Error(t, err)
}
