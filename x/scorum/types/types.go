package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (w SPWithdrawal) CurrentPeriod(t uint64) uint32 {
	return uint32((t - w.CreatedAt) / uint64(w.PeriodDurationInSeconds))
}

func (w SPWithdrawal) PeriodTime(period uint32) uint64 {
	if period > w.TotalPeriods || !w.IsActive {
		return 0
	}
	return w.CreatedAt + uint64(period)*uint64(w.PeriodDurationInSeconds)
}

func (w SPWithdrawal) ToWithdraw(t uint64) sdkmath.Int {
	periodAmount := sdk.NewDecFromInt(w.Total.Int).QuoInt64(int64(w.TotalPeriods)).Ceil().RoundInt()

	withdrew := periodAmount.MulRaw(int64(w.ProcessedPeriod))
	toWithdraw := periodAmount.MulRaw(int64(w.CurrentPeriod(t))).Sub(withdrew)

	if toWithdraw.Add(withdrew).GT(w.Total.Int) {
		toWithdraw = w.Total.Int.Sub(withdrew)
	}

	return toWithdraw
}
