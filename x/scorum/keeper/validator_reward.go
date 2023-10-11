package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

func (k Keeper) PrepareValidatorsReward(ctx sdk.Context) {
	feeCollector := k.accountKeeper.GetModuleAccount(ctx, k.feeCollectorName)

	ctx.Logger().Debug("burn collected gas from fee_collector")
	b := k.bankKeeper.GetBalance(ctx, feeCollector.GetAddress(), types.GasDenom)
	if b.IsPositive() {
		if err := k.Burn(ctx, feeCollector.GetAddress(), b); err != nil {
			panic(fmt.Errorf("failed to burn gas coins: %w", err))
		}
	}

	validatorRewardsParams := k.GetParams(ctx).ValidatorsReward
	if validatorRewardsParams.PoolAddress == "" {
		ctx.Logger().Info("validator rewards pool address is empty")
		ctx.Logger().Debug("skip pouring rewards to fee_collector")

		return
	}

	blockReward := validatorRewardsParams.BlockReward
	if !blockReward.IsValid() || !blockReward.IsPositive() {
		ctx.Logger().Info("validators reward amount is invalid or not positive")
		ctx.Logger().Debug("skip pouring rewards to fee_collector")

		return
	}

	poolAddress := sdk.MustAccAddressFromBech32(validatorRewardsParams.PoolAddress)
	poolBalance := k.bankKeeper.GetBalance(ctx, poolAddress, validatorRewardsParams.BlockReward.Denom)
	if blockReward.IsGTE(poolBalance) {
		blockReward = poolBalance
	}

	if blockReward.IsZero() {
		ctx.Logger().Error("validators reward pool is fully drained")

		return
	}

	if err := k.bankKeeper.SendCoins(ctx, poolAddress, feeCollector.GetAddress(), sdk.NewCoins(blockReward)); err != nil {
		panic(fmt.Errorf("failed to send coins from validators reward pool to fee_collector: %w", err))
	}

	ctx.Logger().Debug("validators reward pool is successfully poured")
}
