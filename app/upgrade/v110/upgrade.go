package v110

import (
	"fmt"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	scorumkeeper "github.com/scorum/cosmos-network/x/scorum/keeper"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const Name = "v110"

func Handler(
	cfg module.Configurator,
	mm *module.Manager,
	cdc *codec.LegacyAmino,
	ps paramtypes.Subspace,
	pk paramskeeper.Keeper,
	cpk *consensusparamkeeper.Keeper,
	sk scorumkeeper.Keeper,
	bk bankkeeper.Keeper,
	stk *stakingkeeper.Keeper,
) func(ctx sdk.Context, _ upgrade.Plan, _ module.VersionMap) (module.VersionMap, error) {
	return func(ctx sdk.Context, _ upgrade.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		baseAppLegacySS := pk.Subspace(baseapp.Paramspace).
			WithKeyTable(paramstypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(ctx, baseAppLegacySS, cpk)

		if err := migrateScorumParams(ctx, cdc, ps); err != nil {
			return nil, fmt.Errorf("fialed to migrate scorum params: %w", err)
		}

		if err := convertSP(ctx, bk, sk, stk); err != nil {
			return nil, fmt.Errorf("fialed to convert sp denom to scr: %w", err)
		}

		return mm.RunMigrations(ctx, cfg, fromVM)
	}
}
