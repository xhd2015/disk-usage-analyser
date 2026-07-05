//go:build darwin

package analyse

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type pnpmSharedContext struct {
	storeIndex map[cloneGroupKey]struct{}
	once       sync.Once
}

func newPnpmSharedContext() *pnpmSharedContext {
	return &pnpmSharedContext{}
}

func (c *pnpmSharedContext) pnpmSharedForNodeModules(nodeModulesPath string) int64 {
	if c == nil {
		return 0
	}
	c.once.Do(func() {
		c.storeIndex = buildPnpmStoreIndex()
	})
	if len(c.storeIndex) == 0 {
		return 0
	}
	return pnpmSharedBytesInNodeModules(nodeModulesPath, c.storeIndex)
}

func resolvePnpmStoreFilesDir() (string, bool) {
	if env := os.Getenv("DISK_USAGE_ANALYSER_PNPM_STORE"); env != "" {
		return env, true
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}

	storeRoot := filepath.Join(home, "Library", "pnpm", "store")
	entries, err := os.ReadDir(storeRoot)
	if err != nil {
		return "", false
	}

	maxV := -1
	var filesDir string
	for _, ent := range entries {
		if !ent.IsDir() {
			continue
		}
		name := ent.Name()
		if !strings.HasPrefix(name, "v") || len(name) < 2 {
			continue
		}
		n, err := strconv.Atoi(name[1:])
		if err != nil || n <= maxV {
			continue
		}
		candidate := filepath.Join(storeRoot, name, "files")
		info, err := os.Stat(candidate)
		if err != nil || !info.IsDir() {
			continue
		}
		maxV = n
		filesDir = candidate
	}
	if maxV < 0 {
		return "", false
	}
	return filesDir, true
}

func buildPnpmStoreIndex() map[cloneGroupKey]struct{} {
	storeDir, ok := resolvePnpmStoreFilesDir()
	if !ok {
		return nil
	}
	return buildCloneStoreIndex(storeDir)
}

func pnpmSharedBytesInNodeModules(nodeModulesPath string, storeIndex map[cloneGroupKey]struct{}) int64 {
	return sharedBytesInNodeModules(nodeModulesPath, storeIndex)
}