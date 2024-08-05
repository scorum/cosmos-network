package v110

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	scorumkeeper "github.com/scorum/cosmos-network/x/scorum/keeper"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
)

func convertSP(
	ctx sdk.Context,
	bk bankkeeper.Keeper, sk scorumkeeper.Keeper,
	gk govkeeper.Keeper, mk mintkeeper.Keeper,
	stk *stakingkeeper.Keeper,
) error {
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
		return fmt.Errorf("failed to set mint params: %w", err)
	}

	govParams := gk.GetParams(ctx)
	for i := range govParams.MinDeposit {
		coin := govParams.MinDeposit[i]
		if coin.Denom == "nsp" {
			coin.Denom = scorumtypes.SCRDenom
		}
		govParams.MinDeposit[i] = coin
	}
	if err := gk.SetParams(ctx, govParams); err != nil {
		return fmt.Errorf("failed to set gov params: %w", err)
	}

	mintParams := mk.GetParams(ctx)
	mintParams.MintDenom = scorumtypes.SCRDenom
	if err := mk.SetParams(ctx, mintParams); err != nil {
		return fmt.Errorf("failed to set mint params: %w", err)
	}

	return nil
}
