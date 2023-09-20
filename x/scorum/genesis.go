package scorum

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scorum/cosmos-network/x/scorum/keeper"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)
	for _, w := range genState.SpWithdrawals {
		k.SetSPWithdrawal(ctx, w)
	}
	for _, v := range genState.RestoreGasAddresses {
		k.SetAddressToRestoreGas(ctx, sdk.MustAccAddressFromBech32(v))
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.SpWithdrawals = k.ListAllWithdrawals(ctx)

	for _, v := range k.ListAddressesForGasRestore(ctx) {
		genesis.RestoreGasAddresses = append(genesis.RestoreGasAddresses, v.String())
	}

	return genesis
}
