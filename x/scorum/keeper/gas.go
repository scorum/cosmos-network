package keeper

import (
	"slices"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

var gasConsumedAddressesPrefix = []byte("gas_consumed")

func (k Keeper) SetAddressToRestoreGas(ctx sdk.Context, addr sdk.AccAddress) {
	s := prefix.NewStore(ctx.KVStore(k.storeKey), gasConsumedAddressesPrefix)

	s.Set(addr.Bytes(), []byte{})
}

func (k Keeper) RestoreGasForAddress(ctx sdk.Context, addr sdk.AccAddress, avgStakedBalance math.LegacyDec, params types.Params) {
	s := prefix.NewStore(ctx.KVStore(k.storeKey), gasConsumedAddressesPrefix)

	gasBalance := k.bankKeeper.GetBalance(ctx, addr, types.GasDenom).Amount
	stakedBalance := math.NewInt(0)
	delegations, err := k.stakingKeeper.GetAllDelegatorDelegations(ctx, addr)
	if err != nil {
		panic(err)
	}
	for _, delegation := range delegations {
		stakedBalance = stakedBalance.Add(delegation.Shares.RoundInt())
	}

	if gasBalance.IsNil() {
		gasBalance = math.ZeroInt()
	}
	if stakedBalance.IsNil() {
		stakedBalance = math.ZeroInt()
	}

	if gasBalance.GTE(params.GasLimit) {
		s.Delete(addr)
	}

	gasAdjust := calculateGasAdjustAmount(
		math.LegacyNewDecFromInt(stakedBalance),
		math.LegacyNewDecFromInt(params.GasLimit),
		math.LegacyNewDecFromInt(params.GasUnconditionedAmount),
		avgStakedBalance,
		params.GasAdjustCoefficient,
	).RoundInt()

	// do not overflow gasLimit
	if gasBalance.Add(gasAdjust).GT(params.GasLimit) {
		gasAdjust = params.GasLimit.Sub(gasBalance)
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

func calculateGasAdjustAmount(stakedBalance, gasLimit, gasUnconditionedAmount, avgStakedBalance, gasAdjustCoefficient math.LegacyDec) math.LegacyDec {
	//                                          stakedBalance
	// adjustAmount = gasUnconditionedAmount + ------------------ * GasLimit * GasAdjustCoefficient
	//                                          avgStakedBalance
	return gasUnconditionedAmount.Add(stakedBalance.Quo(avgStakedBalance).Mul(gasLimit).Mul(gasAdjustCoefficient))
}

func (k Keeper) GetAverageStakedBalance(ctx sdk.Context) math.LegacyDec {
	supervisors := k.GetParams(ctx).Supervisors
	total, size := math.LegacyZeroDec(), int64(0)

	delegations, err := k.stakingKeeper.GetAllDelegations(ctx)
	if err != nil {
		panic(err)
	}
	for _, delegation := range delegations {
		if slices.Contains(supervisors, delegation.DelegatorAddress) {
			continue
		}

		total = total.Add(delegation.Shares)
		size++
	}

	if size == 0 {
		return math.LegacyNewDec(1)
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
