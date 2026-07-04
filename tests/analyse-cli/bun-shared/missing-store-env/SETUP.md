# Scenario

**Leaf**: clone fixture with no env override and no default cache under HOME → `bun_shared=0`

## Steps

1. Create fixture cache file and clone into `node_modules/pkg/file` via `cp -c` (darwin).
2. Do **not** set `DISK_USAGE_ANALYSER_BUN_CACHE`.
3. Point `HOME` at an empty directory (no `.bun/install/cache`).
4. Run `analyse.Analyse` on fixture root.

```go
import (
	"path/filepath"
	"runtime"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if runtime.GOOS != "darwin" {
		t.Skip("missing-store-env clone fixture requires darwin")
	}
	tempRoot := filepath.Dir(req.FixtureDir)
	cacheRoot := bunCacheRoot(tempRoot)
	cacheFile := writeBunCacheFile(t, cacheRoot, file4K)
	cloneCacheFileTo(t, cacheFile, "node_modules/pkg/file", req.FixtureDir)
	t.Setenv("DISK_USAGE_ANALYSER_BUN_CACHE", "")
	isolateHomeWithoutBunCache(t)
	return nil
}
```