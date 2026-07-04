# Scenario

**Leaf**: regular `node_modules` file with no cache clone match yields `bun_shared=0`

## Steps

1. Create fixture cache `pkg@1.0.0@@@1/file` (4096 B) and set env (unmatched content).
2. Write a separate 4096 B regular file at `node_modules/file` (not cloned).
3. Run `analyse.Analyse` on fixture root.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	tempRoot := filepath.Dir(req.FixtureDir)
	cacheRoot := bunCacheRoot(tempRoot)
	writeBunCacheFile(t, cacheRoot, file4K)
	writeSizedFile(t, req.FixtureDir, "node_modules/file", file4K)
	setBunCacheEnv(t, cacheRoot)
	return nil
}
```