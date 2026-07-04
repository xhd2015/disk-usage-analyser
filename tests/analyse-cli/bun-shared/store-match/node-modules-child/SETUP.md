# Scenario

**Leaf**: immediate-child `node_modules` row reports cache-matched `bun_shared`

## Steps

1. Create fixture cache `pkg@1.0.0@@@1/file` (4096 B).
2. Clone cache file into `node_modules/pkg/file` via `cp -c`.
3. Set `DISK_USAGE_ANALYSER_BUN_CACHE` to fixture cache root.
4. Run `analyse.Analyse` on fixture root.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	tempRoot := filepath.Dir(req.FixtureDir)
	cacheRoot := bunCacheRoot(tempRoot)
	cacheFile := writeBunCacheFile(t, cacheRoot, file4K)
	cloneCacheFileTo(t, cacheFile, "node_modules/pkg/file", req.FixtureDir)
	setBunCacheEnv(t, cacheRoot)
	return nil
}
```