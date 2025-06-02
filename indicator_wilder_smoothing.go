package techan

import (
	"github.com/sdcoffey/big"
)

type wilderSmoothingIndicator struct {
	indicator   Indicator
	window      int
	resultCache resultCache
}

func NewWilderSmoothing(indicator Indicator, window int) Indicator {
	return &wilderSmoothingIndicator{
		indicator:   indicator,
		window:      window,
		resultCache: make([]*big.Decimal, 1000),
	}
}

func (ws wilderSmoothingIndicator) cache() resultCache { return ws.resultCache }
func (ws *wilderSmoothingIndicator) setCache(newCache resultCache) {
	ws.resultCache = newCache
}
func (ws wilderSmoothingIndicator) windowSize() int { return ws.window }

func (ws *wilderSmoothingIndicator) Calculate(index int) big.Decimal {
	if cachedValue := returnIfCached(ws, index, nil); cachedValue != nil {
		return *cachedValue
	}

	var result big.Decimal
	if index < ws.window {
		result = ws.indicator.Calculate(index)
	} else if index == ws.window {
		// First value is simple average
		sum := big.ZERO
		for i := 0; i < ws.window; i++ {
			sum = sum.Add(ws.indicator.Calculate(index - i))
		}
		result = sum.Div(big.NewDecimal(float64(ws.window)))
	} else {
		// Subsequent values use Wilder's smoothing formula:
		// (Prior Value * (period - 1) + Current Value) / period
		prevValue := ws.Calculate(index - 1)
		currentValue := ws.indicator.Calculate(index)
		multiplier := big.NewDecimal(float64(ws.window - 1))
		result = prevValue.Mul(multiplier).Add(currentValue).Div(big.NewDecimal(float64(ws.window)))
	}

	cacheResult(ws, index, result)
	return result
} 