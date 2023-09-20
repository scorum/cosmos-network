package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scorum/cosmos-network/testutil/sample"
	"github.com/scorum/cosmos-network/x/scorum/keeper"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServer_StopWithdrawalSp(t *testing.T) {
	set, ctx := setupKeeper(t)

	s := keeper.NewMsgServer(set.keeper)

	addr := sample.AccAddress()
	require.NoError(t, set.keeper.Mint(ctx.Context, addr, sdk.NewCoin(types.SPDenom, sdk.NewInt(100))))

	resp, err := s.WithdrawSP(ctx, &types.MsgWithdrawSP{
		Owner:     addr.String(),
		Recipient: addr.String(),
		Amount:    sdk.IntProto{Int: sdk.NewInt(50)},
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.WithdrawalId)

	_, err = s.StopSPWithdrawal(ctx, &types.MsgStopSPWithdrawal{
		Owner: addr.String(),
		Id:    resp.WithdrawalId,
	})
	require.NoError(t, err)

	act, has := set.keeper.GetSPWithdrawal(ctx.Context, addr, resp.WithdrawalId)
	require.True(t, has)
	require.False(t, act.IsActive)

	_, err = s.StopSPWithdrawal(ctx, &types.MsgStopSPWithdrawal{
		Owner: addr.String(),
		Id:    resp.WithdrawalId,
	})
	require.Error(t, err)
}
