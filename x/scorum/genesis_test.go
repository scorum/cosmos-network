package scorum_test

import (
	"testing"

	keepertest "github.com/scorum/cosmos-network/testutil/keeper"
	"github.com/scorum/cosmos-network/x/scorum"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/stretchr/testify/require"
)

func TestInitGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	ctx := keepertest.GetContext(t)

	k := keepertest.ScorumKeeper(t, ctx)
	scorum.InitGenesis(ctx.Context, k, genesisState)
	require.Equal(t, &types.GenesisState{
		Params:              types.DefaultParams(),
		SpWithdrawals:       []types.SPWithdrawal{},
		RestoreGasAddresses: ([]string)(nil),
	}, scorum.ExportGenesis(ctx.Context, k))
}
