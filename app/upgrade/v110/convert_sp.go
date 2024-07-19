package v110

import (
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	scorumkeeper "github.com/scorum/cosmos-network/x/scorum/keeper"
)

func convertSP(ctx sdk.Context, bk bankkeeper.Keeper, sk scorumkeeper.Keeper, ps paramtypes.Subspace) error {
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

	var stakingParams stakingtypes.Params
	ps.GetParamSet(ctx, &stakingParams)
	stakingParams.BondDenom = scorumtypes.SCRDenom
	ps.SetParamSet(ctx, &stakingParams)

	return nil
}
