package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/cosmos/cosmos-sdk/x/nft/keeper"
)

func NftKeeper(t testing.TB, ctx TestContext) keeper.Keeper {
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	nft.RegisterInterfaces(registry)

	k := keeper.NewKeeper(ctx.KVKeys[nft.StoreKey], cdc, AccountKeeper(t, ctx), BankKeeper(t, ctx))

	return k
}
