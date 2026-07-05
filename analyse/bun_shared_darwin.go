//go:build darwin

package analyse

import (
	"os"
	"path/filepath"
	"sync"
)

type bunSharedContext struct {
	cacheIndex map[cloneGroupKey]struct{}
	once       sync.Once
}

func newBunSharedContext() *bunSharedContext {
	return &bunSharedContext{}
}

func (c *bunSharedContext) bunSharedForNodeModules(nodeModulesPath string) int64 {
	if c == nil {
		return 0
	}
	c.once.Do(func() {
		c.cacheIndex = buildBunCacheIndex()
	})
	if len(c.cacheIndex) == 0 {
		return 0
	}
	return sharedBytesInNodeModules(nodeModulesPath, c.cacheIndex)
}

func resolveBunCacheDir() (string, bool) {
	if env := os.Getenv("DISK_USAGE_ANALYSER_BUN_CACHE"); env != "" {
		return env, true
	}
	if env := os.Getenv("BUN_INSTALL_CACHE_DIR"); env != "" {
		if info, err := os.Stat(env); err == nil && info.IsDir() {
			return env, true
		}
		return "", false
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}
	cacheDir := filepath.Join(home, ".bun", "install", "cache")
	info, err := os.Stat(cacheDir)
	if err != nil || !info.IsDir() {
		return "", false
	}
	return cacheDir, true
}

func buildBunCacheIndex() map[cloneGroupKey]struct{} {
	cacheDir, ok := resolveBunCacheDir()
	if !ok {
		return nil
	}
	return buildCloneStoreIndex(cacheDir)
}