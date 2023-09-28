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

func (w SPWithdrawal) WithdrownByPeriod(p uint32) sdkmath.Int {
	periodAmount := w.Total.Int.QuoRaw(int64(w.TotalPeriods))
	mod := w.Total.Int.ModRaw(int64(w.TotalPeriods))

	return periodAmount.MulRaw(int64(p)).Add(sdk.MinInt(mod, sdk.NewInt(int64(p))))
}

func (w SPWithdrawal) ToWithdraw(t uint64) sdkmath.Int {
	withdrew := w.WithdrownByPeriod(w.ProcessedPeriod)
	toWithdraw := w.WithdrownByPeriod(w.CurrentPeriod(t)).Sub(withdrew)

	if toWithdraw.Add(withdrew).GT(w.Total.Int) {
		toWithdraw = w.Total.Int.Sub(withdrew)
	}

	return toWithdraw
}
