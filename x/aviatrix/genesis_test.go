package aviatrix_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/scorum/cosmos-network/testutil/keeper"
	"github.com/scorum/cosmos-network/x/aviatrix"
	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	ctx := keepertest.GetContext(t)

	k := keepertest.AviatrixKeeper(t, ctx)
	aviatrix.InitGenesis(ctx.Context, k, genesisState)
	require.Equal(t, &types.GenesisState{
		Params: types.DefaultParams(),
	}, aviatrix.ExportGenesis(ctx.Context, k))
}
