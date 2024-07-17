package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/syndtr/goleveldb/leveldb"

	"github.com/syndtr/goleveldb/leveldb/util"

	db "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/store/iavl"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/cosmos/cosmos-sdk/store/types"
	"github.com/scorum/cosmos-network/app"
	"github.com/spf13/cobra"
)

// AddPruneCmd prunes all states except the latest.
func AddPruneCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prune",
		Short: "Prune all states except the latest. Works only with levelDB.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.GetClientContextFromCmd(cmd)
			if ctx.Output == nil {
				ctx = ctx.WithOutput(os.Stdout)
			}
			logger := log.NewTMLogger(ctx.Output)

			logger.Info(fmt.Sprintf("Home dir is %s", ctx.HomeDir))

			if err := prune(ctx, logger); err != nil {
				return fmt.Errorf("failed to prune: %w", err)
			}

			if err := compact(ctx, logger); err != nil {
				return fmt.Errorf("failed to compact: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")

	return cmd
}

func prune(ctx client.Context, logger log.Logger) error {
	db, err := db.NewDB("application", db.GoLevelDBBackend, path.Join(ctx.HomeDir, "data"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	s := rootmulti.NewStore(db, logger)
	s.SetPruning(pruningtypes.NewPruningOptions(pruningtypes.PruningEverything))

	keys := app.GetKeys()
	logger.Info(fmt.Sprintf("%d modules found", len(keys)))

	for _, k := range keys {
		s.MountStoreWithDB(k, types.StoreTypeIAVL, nil)
	}

	if err := s.LoadLatestVersion(); err != nil {
		return fmt.Errorf("failed to load latest version: %w", err)
	}

	latestVersion := s.LatestVersion()
	logger.Info(fmt.Sprintf("Current version: %d", latestVersion))

	for module, key := range keys {
		logger := logger.With("module", module)

		ms := s.GetKVStore(key).(*iavl.Store)
		var toPrune []int64
		for _, v := range ms.GetAllVersions() {
			if int64(v) == latestVersion {
				continue
			}
			toPrune = append(toPrune, int64(v))
		}

		logger.Info(fmt.Sprintf("%d versions to prune", len(toPrune)))
		if err := ms.DeleteVersions(toPrune...); err != nil {
			return fmt.Errorf("failed to delete versions: %w", err)
		}
		logger.Info("pruned")
	}

	logger.Info("Pruning completed")

	return nil
}

func compact(ctx client.Context, logger log.Logger) error {
	logger.Info("Starting compaction")

	db, err := leveldb.OpenFile(path.Join(ctx.HomeDir, "data", "application.db"), nil)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	err = db.CompactRange(util.Range{Start: nil, Limit: nil})
	if err != nil {
		return err
	}

	logger.Info("Compaction completed successfully")

	return nil
}
