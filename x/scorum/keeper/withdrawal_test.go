package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/scorum/cosmos-network/testutil/sample"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
)

func TestKeeper_GetSPWithdrawal(t *testing.T) {
	set, ctx := setupKeeper(t)

	addr := sample.AccAddress()
	set.bankKeeper.InitGenesis(ctx.Context, &banktypes.GenesisState{
		Params: banktypes.DefaultParams(),
		Balances: []banktypes.Balance{
			{addr.String(), sdk.NewCoins(sdk.NewInt64Coin(scorumtypes.SPDenom, 100))},
		},
	})

	id := uuid.New().String()
	exp := scorumtypes.SPWithdrawal{
		Id:                      id,
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
	act, has := set.keeper.GetSPWithdrawal(ctx.Context, addr, id)
	require.True(t, has)
	require.Equal(t, exp, act)

	_, has = set.keeper.GetSPWithdrawal(ctx.Context, sample.AccAddress(), id)
	require.False(t, has)
}

func TestKeeper_ListWithdrawals(t *testing.T) {
	set, ctx := setupKeeper(t)

	addr := sample.AccAddress()
	set.bankKeeper.InitGenesis(ctx.Context, &banktypes.GenesisState{
		Params: banktypes.DefaultParams(),
		Balances: []banktypes.Balance{
			{addr.String(), sdk.NewCoins(sdk.NewInt64Coin(scorumtypes.SPDenom, 100))},
		},
	})

	exp := []scorumtypes.SPWithdrawal{
		{
			Id:                      uuid.New().String(),
			From:                    addr.String(),
			To:                      addr.String(),
			Total:                   sdk.IntProto{Int: sdk.NewInt(50)},
			PeriodDurationInSeconds: 5,
			TotalPeriods:            10,
			ProcessedPeriod:         1,
			IsActive:                false,
			CreatedAt:               0,
		},
		{
			Id:                      uuid.New().String(),
			From:                    addr.String(),
			To:                      addr.String(),
			Total:                   sdk.IntProto{Int: sdk.NewInt(50)},
			PeriodDurationInSeconds: 5,
			TotalPeriods:            10,
			ProcessedPeriod:         1,
			IsActive:                true,
			CreatedAt:               5,
		},
	}
	for _, v := range exp {
		set.keeper.SetSPWithdrawal(ctx.Context, v)
	}

	require.ElementsMatch(t, exp, set.keeper.ListWithdrawals(ctx.Context, addr))
	require.Empty(t, set.keeper.ListWithdrawals(ctx.Context, sample.AccAddress()))
}

func TestKeeper_ListAllWithdrawals(t *testing.T) {
	set, ctx := setupKeeper(t)

	addr1, addr2 := sample.AccAddress(), sample.AccAddress()
	set.bankKeeper.InitGenesis(ctx.Context, &banktypes.GenesisState{
		Params: banktypes.DefaultParams(),
		Balances: []banktypes.Balance{
			{addr1.String(), sdk.NewCoins(sdk.NewInt64Coin(scorumtypes.SPDenom, 100))},
			{addr2.String(), sdk.NewCoins(sdk.NewInt64Coin(scorumtypes.SPDenom, 100))},
		},
	})

	exp := []scorumtypes.SPWithdrawal{
		{
			Id:                      uuid.New().String(),
			From:                    addr1.String(),
			To:                      addr1.String(),
			Total:                   sdk.IntProto{Int: sdk.NewInt(50)},
			PeriodDurationInSeconds: 5,
			TotalPeriods:            10,
			ProcessedPeriod:         1,
			IsActive:                false,
			CreatedAt:               0,
		},
		{
			Id:                      uuid.New().String(),
			From:                    addr2.String(),
			To:                      addr2.String(),
			Total:                   sdk.IntProto{Int: sdk.NewInt(50)},
			PeriodDurationInSeconds: 5,
			TotalPeriods:            10,
			ProcessedPeriod:         1,
			IsActive:                true,
			CreatedAt:               5,
		},
	}
	for _, v := range exp {
		set.keeper.SetSPWithdrawal(ctx.Context, v)
	}

	require.ElementsMatch(t, exp, set.keeper.ListAllWithdrawals(ctx.Context))
}

func TestKeeper_WithdrawSP(t *testing.T) {
	set, ctx := setupKeeper(t)

	addr1, addr2, addr3, addr4 := sample.AccAddress(), sample.AccAddress(), sample.AccAddress(), sample.AccAddress()
	id1, id2, id3, id4 := uuid.New().String(), uuid.New().String(), uuid.New().String(), uuid.New().String()

	set.bankKeeper.InitGenesis(ctx.Context, &banktypes.GenesisState{
		Params: banktypes.DefaultParams(),
		Balances: []banktypes.Balance{
			{addr1.String(), sdk.NewCoins(sdk.NewInt64Coin(scorumtypes.SPDenom, 10))},
			{addr2.String(), sdk.NewCoins(sdk.NewInt64Coin(scorumtypes.SPDenom, 110))},
			{addr3.String(), sdk.NewCoins(sdk.NewInt64Coin(scorumtypes.SPDenom, 1100))},
			{addr4.String(), sdk.NewCoins(sdk.NewInt64Coin(scorumtypes.SPDenom, 1000))},
		},
	})

	set.keeper.SetSPWithdrawal(ctx.Context, scorumtypes.SPWithdrawal{
		Id:                      id1,
		From:                    addr1.String(),
		To:                      addr1.String(),
		Total:                   sdk.IntProto{Int: sdk.NewInt(500)},
		PeriodDurationInSeconds: 5,
		TotalPeriods:            10,
		ProcessedPeriod:         0,
		IsActive:                true,
		CreatedAt:               0,
	})
	set.keeper.SetSPWithdrawal(ctx.Context, scorumtypes.SPWithdrawal{
		Id:                      id2,
		From:                    addr2.String(),
		To:                      addr2.String(),
		Total:                   sdk.IntProto{Int: sdk.NewInt(500)},
		PeriodDurationInSeconds: 5,
		TotalPeriods:            10,
		ProcessedPeriod:         0,
		IsActive:                true,
		CreatedAt:               0,
	})
	set.keeper.SetSPWithdrawal(ctx.Context, scorumtypes.SPWithdrawal{
		Id:                      id3,
		From:                    addr3.String(),
		To:                      addr3.String(),
		Total:                   sdk.IntProto{Int: sdk.NewInt(500)},
		PeriodDurationInSeconds: 5,
		TotalPeriods:            10,
		ProcessedPeriod:         0,
		IsActive:                true,
		CreatedAt:               0,
	})
	set.keeper.SetSPWithdrawal(ctx.Context, scorumtypes.SPWithdrawal{
		Id:                      uuid.NewString(),
		From:                    addr3.String(),
		To:                      addr3.String(),
		Total:                   sdk.IntProto{Int: sdk.NewInt(500)},
		PeriodDurationInSeconds: 5,
		TotalPeriods:            10,
		ProcessedPeriod:         0,
		IsActive:                false,
		CreatedAt:               0,
	})
	set.keeper.SetSPWithdrawal(ctx.Context, scorumtypes.SPWithdrawal{
		Id:                      id4,
		From:                    addr4.String(),
		To:                      addr4.String(),
		Total:                   sdk.IntProto{Int: sdk.NewInt(10)},
		PeriodDurationInSeconds: 5,
		TotalPeriods:            3,
		ProcessedPeriod:         0,
		IsActive:                true,
		CreatedAt:               0,
	})

	set.keeper.WithdrawSP(ctx.Context, 3)
	require.Equal(t, "10nsp", set.bankKeeper.GetBalance(ctx.Context, addr1, scorumtypes.SPDenom).String())
	require.Equal(t, "110nsp", set.bankKeeper.GetBalance(ctx.Context, addr2, scorumtypes.SPDenom).String())
	require.Equal(t, "1100nsp", set.bankKeeper.GetBalance(ctx.Context, addr3, scorumtypes.SPDenom).String())
	require.Equal(t, "1000nsp", set.bankKeeper.GetBalance(ctx.Context, addr4, scorumtypes.SPDenom).String())

	set.keeper.WithdrawSP(ctx.Context, 5)
	require.Equal(t, "0nsp", set.bankKeeper.GetBalance(ctx.Context, addr1, scorumtypes.SPDenom).String())
	require.Equal(t, "10nscr", set.bankKeeper.GetBalance(ctx.Context, addr1, scorumtypes.SCRDenom).String())
	require.Equal(t, "60nsp", set.bankKeeper.GetBalance(ctx.Context, addr2, scorumtypes.SPDenom).String())
	require.Equal(t, "50nscr", set.bankKeeper.GetBalance(ctx.Context, addr2, scorumtypes.SCRDenom).String())
	require.Equal(t, "1050nsp", set.bankKeeper.GetBalance(ctx.Context, addr3, scorumtypes.SPDenom).String())
	require.Equal(t, "50nscr", set.bankKeeper.GetBalance(ctx.Context, addr3, scorumtypes.SCRDenom).String())
	require.Equal(t, "996nsp", set.bankKeeper.GetBalance(ctx.Context, addr4, scorumtypes.SPDenom).String())
	require.Equal(t, "4nscr", set.bankKeeper.GetBalance(ctx.Context, addr4, scorumtypes.SCRDenom).String())

	set.keeper.WithdrawSP(ctx.Context, 6)
	require.Equal(t, "0nsp", set.bankKeeper.GetBalance(ctx.Context, addr1, scorumtypes.SPDenom).String())
	require.Equal(t, "10nscr", set.bankKeeper.GetBalance(ctx.Context, addr1, scorumtypes.SCRDenom).String())
	require.Equal(t, "60nsp", set.bankKeeper.GetBalance(ctx.Context, addr2, scorumtypes.SPDenom).String())
	require.Equal(t, "50nscr", set.bankKeeper.GetBalance(ctx.Context, addr2, scorumtypes.SCRDenom).String())
	require.Equal(t, "1050nsp", set.bankKeeper.GetBalance(ctx.Context, addr3, scorumtypes.SPDenom).String())
	require.Equal(t, "50nscr", set.bankKeeper.GetBalance(ctx.Context, addr3, scorumtypes.SCRDenom).String())
	require.Equal(t, "996nsp", set.bankKeeper.GetBalance(ctx.Context, addr4, scorumtypes.SPDenom).String())
	require.Equal(t, "4nscr", set.bankKeeper.GetBalance(ctx.Context, addr4, scorumtypes.SCRDenom).String())

	set.keeper.WithdrawSP(ctx.Context, 12)
	require.Equal(t, "0nsp", set.bankKeeper.GetBalance(ctx.Context, addr1, scorumtypes.SPDenom).String())
	require.Equal(t, "10nscr", set.bankKeeper.GetBalance(ctx.Context, addr1, scorumtypes.SCRDenom).String())
	require.Equal(t, "10nsp", set.bankKeeper.GetBalance(ctx.Context, addr2, scorumtypes.SPDenom).String())
	require.Equal(t, "100nscr", set.bankKeeper.GetBalance(ctx.Context, addr2, scorumtypes.SCRDenom).String())
	require.Equal(t, "1000nsp", set.bankKeeper.GetBalance(ctx.Context, addr3, scorumtypes.SPDenom).String())
	require.Equal(t, "100nscr", set.bankKeeper.GetBalance(ctx.Context, addr3, scorumtypes.SCRDenom).String())
	require.Equal(t, "993nsp", set.bankKeeper.GetBalance(ctx.Context, addr4, scorumtypes.SPDenom).String())
	require.Equal(t, "7nscr", set.bankKeeper.GetBalance(ctx.Context, addr4, scorumtypes.SCRDenom).String())

	set.keeper.WithdrawSP(ctx.Context, 20)
	require.Equal(t, "0nsp", set.bankKeeper.GetBalance(ctx.Context, addr1, scorumtypes.SPDenom).String())
	require.Equal(t, "10nscr", set.bankKeeper.GetBalance(ctx.Context, addr1, scorumtypes.SCRDenom).String())
	require.Equal(t, "0nsp", set.bankKeeper.GetBalance(ctx.Context, addr2, scorumtypes.SPDenom).String())
	require.Equal(t, "110nscr", set.bankKeeper.GetBalance(ctx.Context, addr2, scorumtypes.SCRDenom).String())
	require.Equal(t, "900nsp", set.bankKeeper.GetBalance(ctx.Context, addr3, scorumtypes.SPDenom).String())
	require.Equal(t, "200nscr", set.bankKeeper.GetBalance(ctx.Context, addr3, scorumtypes.SCRDenom).String())
	require.Equal(t, "990nsp", set.bankKeeper.GetBalance(ctx.Context, addr4, scorumtypes.SPDenom).String())
	require.Equal(t, "10nscr", set.bankKeeper.GetBalance(ctx.Context, addr4, scorumtypes.SCRDenom).String())

	set.keeper.WithdrawSP(ctx.Context, 20)
	require.Equal(t, "0nsp", set.bankKeeper.GetBalance(ctx.Context, addr1, scorumtypes.SPDenom).String())
	require.Equal(t, "10nscr", set.bankKeeper.GetBalance(ctx.Context, addr1, scorumtypes.SCRDenom).String())
	require.Equal(t, "0nsp", set.bankKeeper.GetBalance(ctx.Context, addr2, scorumtypes.SPDenom).String())
	require.Equal(t, "110nscr", set.bankKeeper.GetBalance(ctx.Context, addr2, scorumtypes.SCRDenom).String())
	require.Equal(t, "900nsp", set.bankKeeper.GetBalance(ctx.Context, addr3, scorumtypes.SPDenom).String())
	require.Equal(t, "200nscr", set.bankKeeper.GetBalance(ctx.Context, addr3, scorumtypes.SCRDenom).String())

	set.keeper.WithdrawSP(ctx.Context, 100)
	require.Equal(t, "0nsp", set.bankKeeper.GetBalance(ctx.Context, addr1, scorumtypes.SPDenom).String())
	require.Equal(t, "10nscr", set.bankKeeper.GetBalance(ctx.Context, addr1, scorumtypes.SCRDenom).String())
	require.Equal(t, "0nsp", set.bankKeeper.GetBalance(ctx.Context, addr2, scorumtypes.SPDenom).String())
	require.Equal(t, "110nscr", set.bankKeeper.GetBalance(ctx.Context, addr2, scorumtypes.SCRDenom).String())
	require.Equal(t, "600nsp", set.bankKeeper.GetBalance(ctx.Context, addr3, scorumtypes.SPDenom).String())
	require.Equal(t, "500nscr", set.bankKeeper.GetBalance(ctx.Context, addr3, scorumtypes.SCRDenom).String())
}
