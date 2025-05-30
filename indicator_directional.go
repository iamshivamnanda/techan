package techan

import (
	"github.com/sdcoffey/big"
)

type directionalIndicator struct {
	dm          Indicator
	tr          Indicator
	resultCache resultCache
}

func NewDirectionalIndicator(dm, tr Indicator) Indicator {
	return &directionalIndicator{
		dm:          dm,
		tr:          tr,
		resultCache: make([]*big.Decimal, 1000),
	}
}

func (di directionalIndicator) cache() resultCache { return di.resultCache }
func (di *directionalIndicator) setCache(newCache resultCache) {
	di.resultCache = newCache
}
func (di directionalIndicator) windowSize() int { return 1 }

func (di *directionalIndicator) Calculate(index int) big.Decimal {
	if cachedValue := returnIfCached(di, index, nil); cachedValue != nil {
		return *cachedValue
	}

	tr := di.tr.Calculate(index)
	var result big.Decimal
	if tr.EQ(big.ZERO) {
		result = big.ZERO
	} else {
		// DI = (Smoothed DM / Smoothed TR) * 100
		result = di.dm.Calculate(index).Div(tr).Mul(big.NewDecimal(100))
	}

	cacheResult(di, index, result)
	return result
} 