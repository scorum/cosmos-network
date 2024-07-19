package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

func StakingKeeper(t testing.TB, ctx TestContext) *keeper.Keeper {
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	types.RegisterInterfaces(registry)

	k := keeper.NewKeeper(
		cdc,
		ctx.KVKeys[types.StoreKey],
		AccountKeeper(t, ctx),
		BankKeeper(t, ctx),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Initialize params
	k.SetParams(ctx.Context, types.DefaultParams())

	return k
}
