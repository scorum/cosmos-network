package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/scorum/cosmos-network/testutil/sample"
	"github.com/scorum/cosmos-network/x/scorum/keeper"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/stretchr/testify/require"
)

func TestQuery_QueryWithdrawals(t *testing.T) {
	set, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx.Context)

	addr := sample.AccAddress()

	exp := types.SPWithdrawal{
		Id:                      uuid.NewString(),
		From:                    addr.String(),
		To:                      addr.String(),
		Total:                   sdk.IntProto{Int: sdk.NewInt(50)},
		PeriodDurationInSeconds: 5,
		TotalPeriods:            10,
		ProcessedPeriod:         1,
		IsActive:                true,
		CreatedAt:               0,
	}

	set.keeper.SetSPWithdrawal(ctx.Context, exp)

	srv := keeper.NewQueryServer(set.keeper)
	resp, err := srv.ListWithdrawals(wctx, &types.QueryWithdrawalsRequest{Owner: addr.String()})
	require.NoError(t, err)
	require.Equal(t, []types.SPWithdrawal{exp}, resp.Withdrawals)
}
