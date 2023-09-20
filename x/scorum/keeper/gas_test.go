package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/scorum/cosmos-network/testutil/sample"
	"github.com/stretchr/testify/require"

	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
)

func TestKeeper_getAverageSPBalance(t *testing.T) {
	set, ctx := setupKeeper(t)

	require.True(t, sdk.NewDec(1).Equal(set.keeper.GetAverageSPBalance(ctx.Context)))

	set.bankKeeper.InitGenesis(ctx.Context, &banktypes.GenesisState{
		Params: banktypes.DefaultParams(),
		Balances: []banktypes.Balance{
			{sample.AccAddress().String(), sdk.NewCoins(sdk.NewInt64Coin(scorumtypes.SPDenom, 1))},
			{sample.AccAddress().String(), sdk.NewCoins(sdk.NewInt64Coin(scorumtypes.SPDenom, 2))},
			{sample.AccAddress().String(), sdk.NewCoins(sdk.NewInt64Coin(scorumtypes.SPDenom, 3))},
			{sample.AccAddress().String(), sdk.NewCoins(sdk.NewInt64Coin(scorumtypes.SPDenom, 4))},
			{sample.AccAddress().String(), sdk.NewCoins(sdk.NewInt64Coin(scorumtypes.SPDenom, 5))},
			{sample.AccAddress().String(), sdk.NewCoins(sdk.NewInt64Coin(scorumtypes.SPDenom, 6))},
		},
	})

	exp := sdk.NewDec(21).QuoInt64(6)
	act := set.keeper.GetAverageSPBalance(ctx.Context)
	require.True(t, exp.Equal(act))
}

func TestKeeper_ListAddressesForGasRestore(t *testing.T) {
	set, ctx := setupKeeper(t)

	exp := []sdk.AccAddress{sample.AccAddress(), sample.AccAddress()}

	set.keeper.SetAddressToRestoreGas(ctx.Context, exp[0])
	set.keeper.SetAddressToRestoreGas(ctx.Context, exp[1])

	require.ElementsMatch(t, exp, set.keeper.ListAddressesForGasRestore(ctx.Context))
}
