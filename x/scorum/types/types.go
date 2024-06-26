package types

import (
	sdkmath "cosmossdk.io/math"
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
	periodAmount := w.Total.QuoRaw(int64(w.TotalPeriods))
	mod := w.Total.ModRaw(int64(w.TotalPeriods))

	return periodAmount.MulRaw(int64(p)).Add(sdkmath.MinInt(mod, sdkmath.NewInt(int64(p))))
}

func (w SPWithdrawal) ToWithdraw(t uint64) sdkmath.Int {
	withdrew := w.WithdrownByPeriod(w.ProcessedPeriod)
	toWithdraw := w.WithdrownByPeriod(w.CurrentPeriod(t)).Sub(withdrew)

	if toWithdraw.Add(withdrew).GT(w.Total) {
		toWithdraw = w.Total.Sub(withdrew)
	}

	return toWithdraw
}
