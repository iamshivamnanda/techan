package techan

import "github.com/sdcoffey/big"

type averageTrueRangeIndicator struct {
	series *TimeSeries
	window int
}

// NewAverageTrueRangeIndicator returns a base indicator that calculates the average true range of the
// underlying over a window
// https://www.investopedia.com/terms/a/atr.asp
func NewAverageTrueRangeIndicator(series *TimeSeries, window int) Indicator {
	return averageTrueRangeIndicator{
		series: series,
		window: window,
	}
}

func (atr averageTrueRangeIndicator) Calculate(index int) big.Decimal {
	// Return zero for insufficient data
	if index < atr.window {
		return big.ZERO
	}

	// For the first calculation (when index == window), use simple average
	if index == atr.window {
		sum := big.ZERO
		for i := index; i > index-atr.window; i-- {
			sum = sum.Add(NewTrueRangeIndicator(atr.series).Calculate(i))
		}
		return sum.Div(big.NewFromInt(atr.window))
	}

	// For subsequent periods, use Wilder's smoothing method
	// ATR = [(Prev ATR × (n-1)) + Current TR] / n
	prevATR := atr.Calculate(index - 1)
	currentTR := NewTrueRangeIndicator(atr.series).Calculate(index)

	windowMinusOne := big.NewFromInt(atr.window - 1)
	return (prevATR.Mul(windowMinusOne).Add(currentTR)).Div(big.NewFromInt(atr.window))
}
