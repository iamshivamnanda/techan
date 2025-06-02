package techan

import (
	"github.com/sdcoffey/big"
)

type directionalMovementIndicator struct {
	series      *TimeSeries
	isPlus      bool
	resultCache resultCache
}

func NewDirectionalMovementIndicator(series *TimeSeries, isPlus bool) Indicator {
	return &directionalMovementIndicator{
		series:      series,
		isPlus:      isPlus,
		resultCache: make([]*big.Decimal, 1000),
	}
}

func (dm directionalMovementIndicator) cache() resultCache { return dm.resultCache }
func (dm *directionalMovementIndicator) setCache(newCache resultCache) {
	dm.resultCache = newCache
}
func (dm directionalMovementIndicator) windowSize() int { return 1 }

func (dm *directionalMovementIndicator) Calculate(index int) big.Decimal {
	if cachedValue := returnIfCached(dm, index, nil); cachedValue != nil {
		return *cachedValue
	}

	var result big.Decimal
	if index <= 0 {
		result = big.ZERO
	} else {
		curr := dm.series.Candles[index]
		prev := dm.series.Candles[index-1]

		if dm.isPlus {
			// +DM = Current High - Previous High (if positive and greater than -DM)
			upMove := curr.MaxPrice.Sub(prev.MaxPrice)
			downMove := prev.MinPrice.Sub(curr.MinPrice)
			if upMove.GT(downMove) && upMove.GT(big.ZERO) {
				result = upMove
			} else {
				result = big.ZERO
			}
		} else {
			// -DM = Previous Low - Current Low (if positive and greater than +DM)
			upMove := curr.MaxPrice.Sub(prev.MaxPrice)
			downMove := prev.MinPrice.Sub(curr.MinPrice)
			if downMove.GT(upMove) && downMove.GT(big.ZERO) {
				result = downMove
			} else {
				result = big.ZERO
			}
		}
	}

	cacheResult(dm, index, result)
	return result
} 