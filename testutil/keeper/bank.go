package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/scorum/cosmos-network/app"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
)

func BankKeeper(t testing.TB, ctx TestContext) keeper.Keeper {
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	types.RegisterInterfaces(registry)

	paramsSubspace := typesparams.NewSubspace(cdc,
		codec.NewLegacyAmino(),
		ctx.KVKeys[typesparams.StoreKey],
		ctx.TKeys[typesparams.TStoreKey],
		types.ModuleName,
	)

	k := keeper.NewBaseKeeper(
		cdc,
		ctx.KVKeys[types.StoreKey],
		AccountKeeper(t, ctx),
		paramsSubspace,
		(&app.App{}).BlockedModuleAccountAddrs(),
	)

	// Initialize params
	k.SetParams(ctx.Context, types.Params{
		SendEnabled: []*types.SendEnabled{
			{Denom: scorumtypes.SPDenom, Enabled: false},
			{Denom: scorumtypes.GasDenom, Enabled: false},
		},
		DefaultSendEnabled: true,
	})

	return k
}
