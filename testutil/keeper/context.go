package keeper

import (
	"testing"

	ibcnfttransfertypes "github.com/bianjieai/nft-transfer/types"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icacontrollertypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibchost "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	aviatrixmoduletypes "github.com/scorum/cosmos-network/x/aviatrix/types"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

type TestContext struct {
	sdk.Context

	KVKeys  map[string]*types.KVStoreKey
	MemKeys map[string]*types.MemoryStoreKey
	TKeys   map[string]*types.TransientStoreKey
}

func GetContext(t testing.TB) TestContext {
	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, authz.ModuleName, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey, govtypes.StoreKey,
		paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey, feegrant.StoreKey, evidencetypes.StoreKey,
		ibctransfertypes.StoreKey, icahosttypes.StoreKey, capabilitytypes.StoreKey, group.StoreKey,
		icacontrollertypes.StoreKey, nft.StoreKey, scorumtypes.StoreKey, aviatrixmoduletypes.StoreKey,
		ibcnfttransfertypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)

	for _, k := range keys {
		stateStore.MountStoreWithDB(k, types.StoreTypeIAVL, db)
	}
	for _, k := range memKeys {
		stateStore.MountStoreWithDB(k, types.StoreTypeMemory, db)
	}
	for _, k := range tkeys {
		stateStore.MountStoreWithDB(k, types.StoreTypeTransient, nil)
	}

	require.NoError(t, stateStore.LoadLatestVersion())

	return TestContext{
		Context: sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger()),
		KVKeys:  keys,
		MemKeys: memKeys,
		TKeys:   tkeys,
	}
}
