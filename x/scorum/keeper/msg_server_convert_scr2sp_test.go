package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scorum/cosmos-network/x/scorum/keeper"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServer_ConvertSCR2SP(t *testing.T) {
	set, ctx := setupKeeper(t)

	coin := sdk.NewCoin("scr", sdk.NewInt(1000))

	s := keeper.NewMsgServer(set.keeper)

	require.NoError(t, set.bankKeeper.MintCoins(ctx.Context, types.ModuleName, sdk.NewCoins(coin)))
	require.NoError(t, set.bankKeeper.SendCoinsFromModuleToAccount(ctx.Context, types.ModuleName, set.supervisor, sdk.NewCoins(coin)))

	require.True(t, coin.Equal(set.bankKeeper.GetBalance(ctx.Context, set.supervisor, "scr")))

	_, err := s.ConvertSCR2SP(ctx, &types.MsgConvertSCR2SP{
		Owner:  set.supervisor.String(),
		Amount: sdk.IntProto{Int: sdk.NewInt(300)},
	})
	require.NoError(t, err)

	require.EqualValues(t, 700, set.bankKeeper.GetBalance(ctx.Context, set.supervisor, "scr").Amount.Int64())
	require.EqualValues(t, 300, set.bankKeeper.GetBalance(ctx.Context, set.supervisor, "sp").Amount.Int64())
}
