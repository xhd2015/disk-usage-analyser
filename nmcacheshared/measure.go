package nmcacheshared

import "disk-usage-analyser/analyse"

// NewCalculator returns a calculator that reuses pnpm/bun store indexes across paths.
func NewCalculator() CacheCalculator {
	return &analyseCalculator{inner: analyse.NewCacheSharedCalculator()}
}

// MeasureShared returns pnpm and bun cache-shared bytes for nodeModulesPath.
func MeasureShared(nodeModulesPath string, calc CacheCalculator) (pnpmBytes, bunBytes int64) {
	if calc == nil {
		calc = NewCalculator()
	}
	return calc.PnpmCacheShared(nodeModulesPath), calc.BunCacheShared(nodeModulesPath)
}