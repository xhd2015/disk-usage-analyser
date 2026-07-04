# Scenario

**Leaf**: immediate-child `node_modules` row reports store-matched `pnpm_shared`

## Steps

1. Create fixture store `files/aa/<hash>` (4096 B).
2. Clone store file into `node_modules/pkg/file` via `cp -c`.
3. Set `DISK_USAGE_ANALYSER_PNPM_STORE` to fixture `files/` dir.
4. Run `analyse.Analyse` on fixture root.

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
	setPnpmStoreEnv(t, storeFiles)
	return nil
}
```