# Scenario

**Bug**: summary `bun_shared` preserved when analyse root is `node_modules`

```
# fixture: cache clone in node_modules/pkg/file; analyse node_modules/ directly
Bun cache -> cp -c -> node_modules/pkg/file -> analyse(node_modules) -> summary bun_shared=4096
```

## Steps

1. Create fixture cache `pkg@1.0.0@@@1/file` (4096 B).
2. Clone cache file into `node_modules/pkg/file` via `cp -c`.
3. Set `DISK_USAGE_ANALYSER_BUN_CACHE` to fixture cache root.
4. Point analyse root at `node_modules` (not parent fixture).
5. Run `analyse.Analyse` on the `node_modules` directory.

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
	cloneCacheFileTo(t, cacheFile, "node_modules/pkg/file", fixtureParent)
	setBunCacheEnv(t, cacheRoot)
	req.FixtureDir = filepath.Join(fixtureParent, "node_modules")
	return nil
}
```