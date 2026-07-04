# Scenario

**Bug**: same cache clone key in outer and nested `node_modules` counts once per project tree

```
# two cp -c clones of one cache file under one outer node_modules tree
cache pkg@1.0.0@@@1/file -> node_modules/pkg/file
cache pkg@1.0.0@@@1/file -> node_modules/.pnpm/dep@1.0.0/node_modules/dep/file2

# analyse project root (fixture root)
analyse project root -> bun_shared=4096 (not 8192); bun_shared <= size
```

## Steps

1. Create fixture cache `pkg@1.0.0@@@1/file` (4096 B).
2. Clone cache file into top-level `node_modules/pkg/file` via `cp -c`.
3. Clone the **same** cache file into `node_modules/.pnpm/dep@1.0.0/node_modules/dep/file2` via `cp -c`.
4. Set `DISK_USAGE_ANALYSER_BUN_CACHE` to fixture cache root.
5. Run `analyse.Analyse` on fixture root (project root).

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
	cloneCacheFileTo(t, cacheFile, "node_modules/.pnpm/dep@1.0.0/node_modules/dep/file2", req.FixtureDir)
	setBunCacheEnv(t, cacheRoot)
	return nil
}
```