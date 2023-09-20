package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scorum/cosmos-network/testutil/sample"
	"github.com/scorum/cosmos-network/x/scorum/keeper"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServer_WithdrawSP(t *testing.T) {
	set, ctx := setupKeeper(t)

	s := keeper.NewMsgServer(set.keeper)

	addr := sample.AccAddress()
	require.NoError(t, set.keeper.Mint(ctx.Context, addr, sdk.NewCoin(types.SPDenom, sdk.NewInt(100))))

	_, err := s.WithdrawSP(ctx, &types.MsgWithdrawSP{
		Owner:     addr.String(),
		Recipient: addr.String(),
		Amount:    sdk.IntProto{Int: sdk.NewInt(99)},
	})
	require.Error(t, err)

	resp, err := s.WithdrawSP(ctx, &types.MsgWithdrawSP{
		Owner:     addr.String(),
		Recipient: addr.String(),
		Amount:    sdk.IntProto{Int: sdk.NewInt(50)},
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.WithdrawalId)

	_, err = s.WithdrawSP(ctx, &types.MsgWithdrawSP{
		Owner:     addr.String(),
		Recipient: addr.String(),
		Amount:    sdk.IntProto{Int: sdk.NewInt(50)},
	})
	require.Error(t, err)

	wID := resp.WithdrawalId
	resp, err = s.WithdrawSP(ctx.WithTxBytes([]byte("123")), &types.MsgWithdrawSP{
		Owner:     addr.String(),
		Recipient: addr.String(),
		Amount:    sdk.IntProto{Int: sdk.NewInt(45)},
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.WithdrawalId)
	require.NotEqual(t, wID, resp.WithdrawalId)
}
