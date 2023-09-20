package aviatrix

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scorum/cosmos-network/x/aviatrix/keeper"
	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	return genesis
}
