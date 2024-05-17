//go:build simulation

package app_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/testutil/sims"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	simulationtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/scorum/cosmos-network/app"
	"github.com/stretchr/testify/require"
)

func init() {
	sims.GetSimulatorFlags()
}

// interBlockCacheOpt returns a BaseApp option function that sets the persistent
// inter-block write-through cache.
func interBlockCacheOpt() func(*baseapp.BaseApp) {
	return baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager())
}

func TestAppStateDeterminism(t *testing.T) {
	config := sims.NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = false
	config.AllInvariants = true
	config.NumBlocks = 250
	config.BlockSize = 100
	config.Commit = true
	config.ChainID = sims.SimAppChainID

	numSeeds := 5
	numTimesToRunPerSeed := 3
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	for i := 0; i < numSeeds; i++ {
		config.Seed = rand.Int63()

		for j := 0; j < numTimesToRunPerSeed; j++ {
			var logger log.Logger
			if sims.FlagVerboseValue {
				logger = log.TestingLogger()
			} else {
				logger = log.NewNopLogger()
			}

			db := dbm.NewMemDB()
			app := app.New(
				logger,
				db,
				nil,
				true,
				map[int64]bool{},
				sims.DefaultNodeHome,
				sims.FlagPeriodValue,
				testutil.TestEncodingConfig(),
				sims.AppOptionsMap{},
				interBlockCacheOpt(),
			)

			fmt.Printf(
				"running non-determinism simulation; seed %d: %d/%d, attempt: %d/%d\n",
				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
			)

			_, _, err := simulation.SimulateFromSeed(
				t,
				os.Stdout,
				app.BaseApp,
				sims.AppStateFn(app.AppCodec(), app.SimulationManager()),
				simulationtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
				sims.SimulationOperations(app, app.AppCodec(), config),
				app.ModuleAccountAddrs(),
				config,
				app.AppCodec(),
			)
			require.NoError(t, err)

			if config.Commit {
				PrintStats(db)
			}

			appHash := app.LastCommitID().Hash
			appHashList[j] = appHash

			if j != 0 {
				require.Equal(
					t, string(appHashList[0]), string(appHashList[j]),
					"non-determinism in seed %d: %d/%d, attempt: %d/%d\n", config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
				)
			}
		}
	}
}

// PrintStats prints the corresponding statistics from the app DB.
func PrintStats(db dbm.DB) {
	fmt.Println("\nLevelDB Stats")
	fmt.Println(db.Stats()["leveldb.stats"])
	fmt.Println("LevelDB cached block size", db.Stats()["leveldb.cachedblock"])
}
