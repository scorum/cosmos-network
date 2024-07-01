package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/x/nft"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
)

func AccountKeeper(t testing.TB, ctx TestContext) keeper.AccountKeeper {
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	types.RegisterInterfaces(registry)

	k := keeper.NewAccountKeeper(
		cdc,
		runtime.NewKVStoreService(ctx.KVKeys[types.StoreKey]),
		types.ProtoBaseAccount,
		map[string][]string{
			scorumtypes.ModuleName: {types.Minter, types.Burner},
			nft.ModuleName:         nil,
		},
		address.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Initialize params
	require.NoError(t, k.Params.Set(ctx.Context, types.DefaultParams()))

	return k
}
