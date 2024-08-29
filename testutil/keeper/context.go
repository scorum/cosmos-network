package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	"cosmossdk.io/store/types"
	storetypes "cosmossdk.io/store/types"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	"cosmossdk.io/x/nft"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	db "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	aviatrixtypes "github.com/scorum/cosmos-network/x/aviatrix/types"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/stretchr/testify/require"
)

type TestContext struct {
	sdk.Context

	KVKeys  map[string]*types.KVStoreKey
	MemKeys map[string]*types.MemoryStoreKey
	TKeys   map[string]*types.TransientStoreKey
}

func GetContext(t testing.TB) TestContext {
	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, authz.ModuleName, banktypes.StoreKey, stakingtypes.StoreKey, crisistypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey, govtypes.StoreKey,
		paramstypes.StoreKey, upgradetypes.StoreKey, feegrant.StoreKey, evidencetypes.StoreKey,
		capabilitytypes.StoreKey, group.StoreKey, nft.StoreKey, consensusparamtypes.StoreKey,
		// ibc
		ibcexported.StoreKey, ibcfeetypes.StoreKey,
		icahosttypes.StoreKey, icacontrollertypes.StoreKey,
		// scorum
		scorumtypes.StoreKey, aviatrixtypes.StoreKey,
	)
	tkeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	db := db.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())

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
