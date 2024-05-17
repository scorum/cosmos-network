package v101

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctmmigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"
)

const Name = "v1.0.1 Scorum Upgrade"

func Handler(
	cfg module.Configurator,
	mm *module.Manager,
) func(ctx sdk.Context, _ upgrade.Plan, _ module.VersionMap) (module.VersionMap, error) {
	return func(ctx sdk.Context, _ upgrade.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// prune expired tendermint consensus states to save storage space
		_, err := ibctmmigrations.PruneExpiredConsensusStates(ctx, app.Codec, app.IBCKeeper.ClientKeeper)
		if err != nil {
			return nil, err
		}

		return mm.RunMigrations(ctx, cfg, fromVM)
	}
}
