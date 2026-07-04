# Scenario

**Bug**: same store clone key in outer and nested `node_modules` counts once per project tree

```
# two cp -c clones of one store file under one outer node_modules tree
store files/aa/<hash> -> node_modules/pkg/file
store files/aa/<hash> -> node_modules/.pnpm/dep@1.0.0/node_modules/dep/file2

# analyse project root (fixture root)
analyse project root -> pnpm_shared=4096 (not 8192); pnpm_shared <= size
```

## Steps

1. Create fixture store `files/aa/<hash>` (4096 B).
2. Clone store file into top-level `node_modules/pkg/file` via `cp -c`.
3. Clone the **same** store file into `node_modules/.pnpm/dep@1.0.0/node_modules/dep/file2` via `cp -c`.
4. Set `DISK_USAGE_ANALYSER_PNPM_STORE` to fixture `files/` dir.
5. Run `analyse.Analyse` on fixture root (project root).

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	tempRoot := filepath.Dir(req.FixtureDir)
	storeFiles := pnpmStoreFilesDir(tempRoot)
	storeFile := writePnpmStoreFile(t, storeFiles, file4K)
	cloneStoreFileTo(t, storeFile, "node_modules/pkg/file", req.FixtureDir)
	cloneStoreFileTo(t, storeFile, "node_modules/.pnpm/dep@1.0.0/node_modules/dep/file2", req.FixtureDir)
	setPnpmStoreEnv(t, storeFiles)
	return nil
}
```