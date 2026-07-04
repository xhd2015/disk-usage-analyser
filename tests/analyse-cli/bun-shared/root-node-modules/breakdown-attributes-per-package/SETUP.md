# Scenario

**Bug**: per-package breakdown rows show `bun_shared=0` when analyse root is `node_modules`

Reproduces `disk-usage-analyser analyse node_modules/` on a Bun hoisted install:
summary `bun_shared` is correct (e.g. 172.2M) but every package row in the
breakdown shows `0 B` except nested `node_modules` folders.

```
# fixture: pkg-a has cache clone; pkg-b is regular; analyse node_modules/
cache -> cp -c -> node_modules/pkg-a/file
write regular -> node_modules/pkg-b/file
analyse(node_modules) -> summary bun_shared=4096; pkg-a row should carry 4096
```

## Steps

1. Create fixture Bun cache file (4096 B).
2. Clone into `node_modules/pkg-a/file` via `cp -c`.
3. Write a non-cloned regular file in `node_modules/pkg-b/file`.
4. Set `DISK_USAGE_ANALYSER_BUN_CACHE` to fixture cache root.
5. Analyse the `node_modules` directory (not parent).

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	tempRoot := filepath.Dir(req.FixtureDir)
	fixtureParent := req.FixtureDir
	cacheRoot := bunCacheRoot(tempRoot)
	cacheFile := writeBunCacheFile(t, cacheRoot, file4K)
	cloneCacheFileTo(t, cacheFile, "node_modules/pkg-a/file", fixtureParent)
	writeSizedFile(t, fixtureParent, "node_modules/pkg-b/file", file4K)
	setBunCacheEnv(t, cacheRoot)
	req.FixtureDir = filepath.Join(fixtureParent, "node_modules")
	return nil
}
```