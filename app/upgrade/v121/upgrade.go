package v121

import (
	"context"

	upgrade "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/types/module"
)

const Name = "v1.2.1"

func Handler(
	cfg module.Configurator,
	mm *module.Manager,
) func(ctx context.Context, _ upgrade.Plan, _ module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, _ upgrade.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return mm.RunMigrations(ctx, cfg, fromVM)
	}
}
