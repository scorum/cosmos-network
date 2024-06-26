package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

var gasConsumedAddressesPrefix = []byte("gas_consumed")

func (k Keeper) SetAddressToRestoreGas(ctx sdk.Context, addr sdk.AccAddress) {
	s := prefix.NewStore(ctx.KVStore(k.storeKey), gasConsumedAddressesPrefix)

	s.Set(addr.Bytes(), []byte{})
}

func (k Keeper) RestoreGasForAddress(ctx sdk.Context, addr sdk.AccAddress, avgSPBalance math.LegacyDec, params types.Params) {
	s := prefix.NewStore(ctx.KVStore(k.storeKey), gasConsumedAddressesPrefix)

	gasBalance := k.bankKeeper.GetBalance(ctx, addr, types.GasDenom).Amount
	spBalance := k.bankKeeper.GetBalance(ctx, addr, types.SPDenom).Amount

	if gasBalance.IsNil() {
		gasBalance = math.ZeroInt()
	}
	if spBalance.IsNil() {
		spBalance = math.ZeroInt()
	}

	if gasBalance.GTE(params.GasLimit) {
		s.Delete(addr)
	}

	gasAdjust := calculateGasAdjustAmount(
		math.LegacyNewDecFromInt(spBalance),
		math.LegacyNewDecFromInt(params.GasLimit),
		math.LegacyNewDecFromInt(params.GasUnconditionedAmount),
		avgSPBalance,
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

	avgSPBalance, params := k.GetAverageSPBalance(ctx), k.GetParams(ctx)
	for ; it.Valid(); it.Next() {
		k.RestoreGasForAddress(ctx, it.Key(), avgSPBalance, params)
	}
}

func calculateGasAdjustAmount(spBalance, gasLimit, gasUnconditionedAmount, avgSPBalance, gasAdjustCoefficient math.LegacyDec) math.LegacyDec {
	//                                           spBalance
	// adjustAmount = gasUnconditionedAmount + ------------- * GasLimit * GasAdjustCoefficient
	//                                          avgSPBalance
	return gasUnconditionedAmount.Add(spBalance.Quo(avgSPBalance).Mul(gasLimit).Mul(gasAdjustCoefficient))
}

func (k Keeper) GetAverageSPBalance(ctx sdk.Context) math.LegacyDec {
	supervisors := k.GetParams(ctx).Supervisors
	total, size := math.LegacyZeroDec(), int64(0)
	k.bankKeeper.IterateAllBalances(ctx, func(addr sdk.AccAddress, coin sdk.Coin) (stop bool) {
		if contains(supervisors, addr.String()) {
			return false
		}

		if _, ok := k.accountKeeper.GetAccount(ctx, addr).(authtypes.ModuleAccountI); ok {
			return false
		}

		if coin.Denom == types.SPDenom {
			total = total.Add(math.LegacyNewDecFromInt(coin.Amount))
			size++
		}

		return false
	})

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
