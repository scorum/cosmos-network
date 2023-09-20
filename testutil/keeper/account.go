package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
)

func AccountKeeper(t testing.TB, ctx TestContext) keeper.AccountKeeper {
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	types.RegisterInterfaces(registry)

	paramsSubspace := typesparams.NewSubspace(cdc,
		codec.NewLegacyAmino(),
		ctx.KVKeys[typesparams.StoreKey],
		ctx.TKeys[typesparams.TStoreKey],
		types.ModuleName,
	)

	k := keeper.NewAccountKeeper(
		cdc,
		ctx.KVKeys[types.StoreKey],
		paramsSubspace,
		types.ProtoBaseAccount,
		map[string][]string{
			scorumtypes.ModuleName: {types.Minter, types.Burner},
			nft.ModuleName:         nil,
		},
		"scorum",
	)

	// Initialize params
	k.SetParams(ctx.Context, types.DefaultParams())

	return k
}
