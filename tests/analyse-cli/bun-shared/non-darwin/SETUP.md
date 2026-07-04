# Scenario

**Leaf**: non-darwin platforms always report `bun_shared=0` even with `node_modules` and cache env

## Steps

1. Skip on darwin (APFS cache-match leaves cover darwin).
2. Write `node_modules/file` (4096 B) and set cache env to a fixture cache file.
3. Run `analyse.Analyse` on fixture root.

```go
import (
	"path/filepath"
	"runtime"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if runtime.GOOS == "darwin" {
		t.Skip("non-darwin platform leaf")
	}
	tempRoot := filepath.Dir(req.FixtureDir)
	cacheRoot := bunCacheRoot(tempRoot)
	writeBunCacheFile(t, cacheRoot, file4K)
	writeSizedFile(t, req.FixtureDir, "node_modules/file", file4K)
	setBunCacheEnv(t, cacheRoot)
	return nil
}
```