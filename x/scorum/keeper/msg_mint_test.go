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

func TestMsgServer_Mint(t *testing.T) {
	set, ctx := setupKeeper(t)

	addr := sample.AccAddress()
	coin := sdk.NewCoin(types.SCRDenom, math.NewInt(1000))

	s := keeper.NewMsgServer(set.keeper)

	_, err := s.Mint(ctx.Context, &types.MsgMint{
		Authority: "gov",
		Recipient: addr.String(),
		Amount:    coin,
	})
	require.NoError(t, err)

	require.True(t, coin.Equal(set.bankKeeper.GetBalance(ctx.Context, addr, types.SCRDenom)))
}
