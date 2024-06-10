package keeper

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/scorum/cosmos-network/app"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
)

func BankKeeper(t testing.TB, ctx TestContext) keeper.Keeper {
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	types.RegisterInterfaces(registry)

	k := keeper.NewBaseKeeper(
		cdc,
		ctx.KVKeys[types.StoreKey],
		AccountKeeper(t, ctx),
		(&app.App{}).BlockedModuleAccountAddrs(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
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
