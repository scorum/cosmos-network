package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/runtime"

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
		runtime.NewKVStoreService(ctx.KVKeys[types.StoreKey]),
		AccountKeeper(t, ctx),
		(&app.App{}).BlockedModuleAccountAddrs(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		log.NewNopLogger(),
	)

	// Initialize params
	require.NoError(t, k.SetParams(ctx.Context, types.Params{
		DefaultSendEnabled: true,
	}))
	k.SetSendEnabled(ctx, scorumtypes.GasDenom, false)

	return k
}
