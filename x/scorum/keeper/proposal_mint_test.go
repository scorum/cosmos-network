package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scorum/cosmos-network/testutil/sample"
	"github.com/scorum/cosmos-network/x/scorum/keeper"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/stretchr/testify/require"
)

func TestProposal_Mint(t *testing.T) {
	set, ctx := setupKeeper(t)

	addr := sample.AccAddress()
	coin := sdk.NewCoin("scr", sdk.NewInt(1000))

	require.NoError(t, keeper.HandleMintProposal(ctx.Context, set.keeper, &types.MintProposal{
		Title:       "test",
		Description: "test description",
		Recipient:   addr.String(),
		Amount:      coin,
	}))

	require.True(t, coin.Equal(set.bankKeeper.GetBalance(ctx.Context, addr, "scr")))
}
