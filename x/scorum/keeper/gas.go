package keeper

import (
	"slices"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

var gasConsumedAddressesPrefix = []byte("gas_consumed")

func (k Keeper) SetAddressToRestoreGas(ctx sdk.Context, addr sdk.AccAddress) {
	s := prefix.NewStore(ctx.KVStore(k.storeKey), gasConsumedAddressesPrefix)

	s.Set(addr.Bytes(), []byte{})
}

func (k Keeper) RestoreGasForAddress(ctx sdk.Context, addr sdk.AccAddress, avgStakedBalance sdk.Dec, params types.Params) {
	s := prefix.NewStore(ctx.KVStore(k.storeKey), gasConsumedAddressesPrefix)

	gasBalance := k.bankKeeper.GetBalance(ctx, addr, types.GasDenom).Amount
	stakedBalance := sdk.NewInt(0)
	for _, delegation := range k.stakingKeeper.GetAllDelegatorDelegations(ctx, addr) {
		stakedBalance = stakedBalance.Add(delegation.Shares.RoundInt())
	}

	if gasBalance.IsNil() {
		gasBalance = sdk.ZeroInt()
	}
	if stakedBalance.IsNil() {
		stakedBalance = sdk.ZeroInt()
	}

	if gasBalance.GTE(params.GasLimit.Int) {
		s.Delete(addr)
	}

	gasAdjust := calculateGasAdjustAmount(
		sdk.NewDecFromInt(stakedBalance),
		sdk.NewDecFromInt(params.GasLimit.Int),
		sdk.NewDecFromInt(params.GasUnconditionedAmount.Int),
		avgStakedBalance,
		params.GasAdjustCoefficient.Dec,
	).RoundInt()

	// do not overflow gasLimit
	if gasBalance.Add(gasAdjust).GT(params.GasLimit.Int) {
		gasAdjust = params.GasLimit.Int.Sub(gasBalance)
	}

	if gasAdjust.IsPositive() {
		if err := k.Mint(ctx, addr, sdk.NewCoin(types.GasDenom, gasAdjust)); err != nil {
			panic(err)
		}
	}
}

func (k Keeper) RestoreGas(ctx sdk.Context) {
	s := prefix.NewStore(ctx.KVStore(k.storeKey), gasConsumedAddressesPrefix)
	it := s.Iterator(nil, nil)
	defer it.Close()

	avgStakedBalance, params := k.GetAverageStakedBalance(ctx), k.GetParams(ctx)
	for ; it.Valid(); it.Next() {
		k.RestoreGasForAddress(ctx, it.Key(), avgStakedBalance, params)
	}
}

func calculateGasAdjustAmount(stakedBalance, gasLimit, gasUnconditionedAmount, avgStakedBalance, gasAdjustCoefficient sdk.Dec) sdk.Dec {
	//                                          stakedBalance
	// adjustAmount = gasUnconditionedAmount + ------------------ * GasLimit * GasAdjustCoefficient
	//                                          avgStakedBalance
	return gasUnconditionedAmount.Add(stakedBalance.Quo(avgStakedBalance).Mul(gasLimit).Mul(gasAdjustCoefficient))
}

func (k Keeper) GetAverageStakedBalance(ctx sdk.Context) sdk.Dec {
	supervisors := k.GetParams(ctx).Supervisors
	total, size := sdk.ZeroDec(), int64(0)

	for _, delegation := range k.stakingKeeper.GetAllDelegations(ctx) {
		if slices.Contains(supervisors, delegation.DelegatorAddress) {
			continue
		}

		total = total.Add(delegation.Shares)
		size++
	}

	if size == 0 {
		return sdk.NewDec(1)
	}

	return total.QuoInt64(size)
}

func (k Keeper) ListAddressesForGasRestore(ctx sdk.Context) []sdk.AccAddress {
	s := prefix.NewStore(ctx.KVStore(k.storeKey), gasConsumedAddressesPrefix)
	it := s.Iterator(nil, nil)
	defer it.Close()

	out := make([]sdk.AccAddress, 0)
	for ; it.Valid(); it.Next() {
		out = append(out, it.Key())
	}

	return out
}
