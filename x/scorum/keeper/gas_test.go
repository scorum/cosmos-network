package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scorum/cosmos-network/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestKeeper_ListAddressesForGasRestore(t *testing.T) {
	set, ctx := setupKeeper(t)

	exp := []sdk.AccAddress{sample.AccAddress(), sample.AccAddress()}

	set.keeper.SetAddressToRestoreGas(ctx.Context, exp[0])
	set.keeper.SetAddressToRestoreGas(ctx.Context, exp[1])

	require.ElementsMatch(t, exp, set.keeper.ListAddressesForGasRestore(ctx.Context))
}
