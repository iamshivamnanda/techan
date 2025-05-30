package techan

import (
	"github.com/sdcoffey/big"
)

type averageDirectionalIndexIndicator struct {
	series      *TimeSeries
	window      int
	plusDI      Indicator
	minusDI     Indicator
	resultCache resultCache
}

// NewAverageDirectionalIndexIndicator returns an indicator that calculates the ADX
// which measures trend strength irrespective of direction
// window: number of periods (typically 14)
func NewAverageDirectionalIndexIndicator(series *TimeSeries, window int) Indicator {
	tr := NewTrueRangeIndicator(series)
	plusDM := NewDirectionalMovementIndicator(series, true)
	minusDM := NewDirectionalMovementIndicator(series, false)

	// Smooth TR, +DM, -DM using Wilder's smoothing
	smoothedTR := NewWilderSmoothing(tr, window)
	smoothedPlusDM := NewWilderSmoothing(plusDM, window)
	smoothedMinusDM := NewWilderSmoothing(minusDM, window)

	// Calculate +DI and -DI
	plusDI := NewDirectionalIndicator(smoothedPlusDM, smoothedTR)
	minusDI := NewDirectionalIndicator(smoothedMinusDM, smoothedTR)

	return &averageDirectionalIndexIndicator{
		series:      series,
		window:      window,
		plusDI:      plusDI,
		minusDI:     minusDI,
		resultCache: make([]*big.Decimal, 1000),
	}
}

func (adx *averageDirectionalIndexIndicator) cache() resultCache { return adx.resultCache }
func (adx *averageDirectionalIndexIndicator) setCache(newCache resultCache) {
	adx.resultCache = newCache
}
func (adx averageDirectionalIndexIndicator) windowSize() int { return adx.window }

func (adx *averageDirectionalIndexIndicator) Calculate(index int) big.Decimal {
	if cachedValue := returnIfCached(adx, index, nil); cachedValue != nil {
		return *cachedValue
	}

	var result big.Decimal
	if index < adx.window {
		result = big.ZERO
	} else {
		plusDI := adx.plusDI.Calculate(index)
		minusDI := adx.minusDI.Calculate(index)

		// Calculate DX = |(+DI - -DI)| / (+DI + -DI) * 100
		sumDI := plusDI.Add(minusDI)
		if sumDI.EQ(big.ZERO) {
			result = big.ZERO
		} else {
			diffDI := plusDI.Sub(minusDI).Abs()
			dx := diffDI.Div(sumDI).Mul(big.NewDecimal(100))

			// For first window periods, return DX
			if index == adx.window {
				result = dx
			} else {
				// After window periods, use Wilder's smoothing for ADX
				prevADX := adx.Calculate(index - 1)
				result = prevADX.Mul(big.NewDecimal(float64(adx.window - 1))).Add(dx).Div(big.NewDecimal(float64(adx.window)))
			}
		}
	}

	cacheResult(adx, index, result)
	return result
}

