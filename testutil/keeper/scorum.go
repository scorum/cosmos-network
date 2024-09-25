package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/scorum/cosmos-network/x/scorum/keeper"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

func ScorumKeeper(
	t testing.TB,
	ctx TestContext,
) keeper.Keeper {
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	types.RegisterInterfaces(registry)

	paramsSubspace := typesparams.NewSubspace(cdc,
		codec.NewLegacyAmino(),
		ctx.KVKeys[typesparams.StoreKey],
		ctx.TKeys[typesparams.TStoreKey],
		types.ModuleName,
	)

	k := keeper.NewKeeper(
		cdc,
		ctx.KVKeys[types.StoreKey],
		paramsSubspace,
		AccountKeeper(t, ctx),
		BankKeeper(t, ctx),
		StakingKeeper(t, ctx),
		authtypes.FeeCollectorName,
		"gov",
	)

	// Initialize params
	k.SetParams(ctx.Context, types.DefaultParams())

	return k
}
