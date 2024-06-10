package wrap

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/mikluke/co-pilot/slice"
	scorumkeeper "github.com/scorum/cosmos-network/x/scorum/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

// StakingBankKeeper is a wrapper of cosmos-sdk/x/auth/keeper.BankKeeper that transfers coins to staking reward pool
// instead of burning it. It's used in slashing mechanism to prevent decreasing emission.
type StakingBankKeeper struct {
	bankkeeper.Keeper
	sk *scorumkeeper.Keeper
}

func NewStakingBankKeeper(bk bankkeeper.Keeper, sk *scorumkeeper.Keeper) bankkeeper.Keeper {
	return StakingBankKeeper{
		Keeper: bk,
		sk:     sk,
	}
}

func (k StakingBankKeeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	if !slice.Contains([]string{stakingtypes.BondedPoolName, stakingtypes.NotBondedPoolName}, moduleName) {
		return k.Keeper.BurnCoins(ctx, moduleName, amt)
	}

	pool := k.sk.GetParams(ctx).ValidatorsReward.PoolAddress
	if pool != "" {
		ctx.Logger().Info("send slashed coins to validators reward pool")
		return k.SendCoinsFromModuleToAccount(ctx, moduleName, sdk.MustAccAddressFromBech32(pool), amt)
	}

	ctx.Logger().Error("staking module requested coins burning, but there is no pool to transfer it")
	return k.Keeper.BurnCoins(ctx, moduleName, amt)
}
