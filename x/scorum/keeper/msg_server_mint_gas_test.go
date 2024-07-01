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

func TestMsgServer_MintGas(t *testing.T) {
	set, ctx := setupKeeper(t)

	coin := sdk.NewCoin(types.GasDenom, math.NewInt(1000))

	s := keeper.NewMsgServer(set.keeper)

	require.NoError(t, set.keeper.Mint(ctx.Context, set.supervisor, coin))

	require.True(t, coin.Equal(set.bankKeeper.GetBalance(ctx.Context, set.supervisor, types.GasDenom)))

	_, err := s.MintGas(ctx, &types.MsgMintGas{
		Supervisor: set.supervisor.String(),
		Address:    set.supervisor.String(),
		Amount:     math.NewInt(10000000),
	})
	require.NoError(t, err)

	require.Equal(t, "10001000", set.bankKeeper.GetBalance(ctx.Context, set.supervisor, types.GasDenom).Amount.String())
}

func TestMsgServer_MintGas_NotSupervisor(t *testing.T) {
	set, ctx := setupKeeper(t)

	addr := sample.AccAddress()
	coin := sdk.NewCoin(types.GasDenom, math.NewInt(1000))

	s := keeper.NewMsgServer(set.keeper)

	require.NoError(t, set.keeper.Mint(ctx.Context, addr, coin))

	require.True(t, coin.Equal(set.bankKeeper.GetBalance(ctx.Context, addr, types.GasDenom)))

	_, err := s.MintGas(ctx, &types.MsgMintGas{
		Supervisor: addr.String(),
		Address:    set.supervisor.String(),
		Amount:     math.NewInt(10000000),
	})
	require.Error(t, err)
}
