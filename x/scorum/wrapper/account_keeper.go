package wrap

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/auth/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	scorumkeeper "github.com/scorum/cosmos-network/x/scorum/keeper"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
)

// AccountKeeper is a wrapper of cosmos-sdk/x/auth/keeper.AccountKeeperI that mints gas on setting new account.
// It's used to allow free-gas transactions without registration.
type AccountKeeper struct {
	accountkeeper.AccountKeeper
	bk bankkeeper.Keeper
	sk scorumkeeper.Keeper
}

func NewAccountKeeper(ak accountkeeper.AccountKeeper, bk bankkeeper.Keeper, sk scorumkeeper.Keeper) AccountKeeper {
	return AccountKeeper{
		AccountKeeper: ak,
		bk:            bk,
		sk:            sk,
	}
}

func (k AccountKeeper) SetAccount(ctx sdk.Context, acc types.AccountI) {
	hasAccount := k.AccountKeeper.HasAccount(ctx, acc.GetAddress())
	// must be set before minting to avoid recursion (BankKeeper calls SetAccount if it's not created yet)
	k.AccountKeeper.SetAccount(ctx, acc)
	if !hasAccount {
		if err := k.sk.Mint(ctx, acc.GetAddress(), sdk.NewCoin(scorumtypes.GasDenom, k.sk.GetParams(ctx).GasLimit.Int)); err != nil {
			panic(fmt.Sprintf("failed to mint gas to new account: %s", err.Error()))
		}
	}
}
