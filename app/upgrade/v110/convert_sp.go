package v110

import (
	"fmt"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	scorumkeeper "github.com/scorum/cosmos-network/x/scorum/keeper"
)

func convertSP(ctx sdk.Context, bk bankkeeper.Keeper, sk scorumkeeper.Keeper, stk *stakingkeeper.Keeper) error {
	var err error
	bk.IterateAllBalances(ctx, func(address sdk.AccAddress, coin sdk.Coin) (stop bool) {
		if coin.Denom != "nsp" {
			return false
		}

		if err = sk.Burn(ctx, address, coin); err != nil {
			err = fmt.Errorf("failed to burn nsp")
			return true
		}

		coin.Denom = scorumtypes.SCRDenom
		if err = sk.Mint(ctx, address, coin); err != nil {
			err = fmt.Errorf("failed to mint nscr")
			return true
		}

		return false
	})
	if err != nil {
		return err
	}

	stakingParams := stk.GetParams(ctx)
	stakingParams.BondDenom = scorumtypes.SCRDenom
	if err := stk.SetParams(ctx, stakingParams); err != nil {
		return fmt.Errorf("failed to set params: %w", err)
	}

	return nil
}
