package analyse

import "path/filepath"

// CacheSharedCalculator reuses pnpm store and bun cache indexes across many paths.
type CacheSharedCalculator struct {
	pnpm *pnpmSharedContext
	bun  *bunSharedContext
}

// NewCacheSharedCalculator returns a calculator that builds store indexes lazily on first use.
func NewCacheSharedCalculator() *CacheSharedCalculator {
	return &CacheSharedCalculator{
		pnpm: newPnpmSharedContext(),
		bun:  newBunSharedContext(),
	}
}

// PnpmCacheShared returns bytes in nodeModulesPath that APFS-clone-share with the pnpm store.
func (c *CacheSharedCalculator) PnpmCacheShared(nodeModulesPath string) int64 {
	return c.pnpm.pnpmSharedForNodeModules(absPath(nodeModulesPath))
}

// BunCacheShared returns bytes in nodeModulesPath that APFS-clone-share with the bun install cache.
func (c *CacheSharedCalculator) BunCacheShared(nodeModulesPath string) int64 {
	return c.bun.bunSharedForNodeModules(absPath(nodeModulesPath))
}

func absPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}