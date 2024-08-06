package v110

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	scorumkeeper "github.com/scorum/cosmos-network/x/scorum/keeper"
)

const Name = "v1.1.0"

func Handler(
	cfg module.Configurator,
	mm *module.Manager,
	cdc *codec.LegacyAmino,
	scorumParamSpace paramtypes.Subspace,
	pk paramskeeper.Keeper,
	cpk *consensusparamkeeper.Keeper,
	sk scorumkeeper.Keeper,
	bk bankkeeper.Keeper,
	gk govkeeper.Keeper,
	mk mintkeeper.Keeper,
	stk *stakingkeeper.Keeper,
) func(ctx sdk.Context, _ upgrade.Plan, _ module.VersionMap) (module.VersionMap, error) {
	return func(ctx sdk.Context, _ upgrade.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		baseAppLegacySS := pk.Subspace(baseapp.Paramspace).
			WithKeyTable(paramstypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(ctx, baseAppLegacySS, cpk)

		migrationMap, err := mm.RunMigrations(ctx, cfg, fromVM)
		if err != nil {
			return nil, err
		}

		if err := migrateScorumParams(ctx, cdc, scorumParamSpace); err != nil {
			return nil, fmt.Errorf("fialed to migrate scorum params: %w", err)
		}

		if err := convertSP(ctx, bk, sk, gk, mk, stk); err != nil {
			return nil, fmt.Errorf("fialed to convert sp denom to scr: %w", err)
		}

		return migrationMap, nil
	}
}
