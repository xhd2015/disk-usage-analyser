# Scenario

**Leaf**: regular `node_modules` file with no store clone match yields `pnpm_shared=0`

## Steps

1. Create fixture store `files/aa/<hash>` (4096 B) and set env (unmatched content).
2. Write a separate 4096 B regular file at `node_modules/file` (not cloned).
3. Run `analyse.Analyse` on fixture root.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	tempRoot := filepath.Dir(req.FixtureDir)
	storeFiles := pnpmStoreFilesDir(tempRoot)
	writePnpmStoreFile(t, storeFiles, file4K)
	writeSizedFile(t, req.FixtureDir, "node_modules/file", file4K)
	setPnpmStoreEnv(t, storeFiles)
	return nil
}
```