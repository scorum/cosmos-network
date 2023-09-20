package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/scorum/cosmos-network/x/aviatrix/keeper"
	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

func AviatrixKeeper(t testing.TB, ctx TestContext) keeper.Keeper {
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
		NftKeeper(t, ctx),
		ScorumKeeper(t, ctx),
	)

	// Initialize params
	k.SetParams(ctx.Context, types.DefaultParams())

	return k
}
