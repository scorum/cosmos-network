package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/x/nft"
	"cosmossdk.io/x/nft/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

func NftKeeper(t testing.TB, ctx TestContext) keeper.Keeper {
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	nft.RegisterInterfaces(registry)

	k := keeper.NewKeeper(
		runtime.NewKVStoreService(ctx.KVKeys[nft.StoreKey]),
		cdc,
		AccountKeeper(t, ctx),
		BankKeeper(t, ctx),
	)

	return k
}
