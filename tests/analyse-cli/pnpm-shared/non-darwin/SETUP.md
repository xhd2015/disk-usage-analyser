# Scenario

**Leaf**: non-darwin platforms always report `pnpm_shared=0` even with `node_modules` and store env

## Steps

1. Skip on darwin (APFS store-match leaves cover darwin).
2. Write `node_modules/file` (4096 B) and set store env to a fixture store file.
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
	storeFiles := pnpmStoreFilesDir(tempRoot)
	writePnpmStoreFile(t, storeFiles, file4K)
	writeSizedFile(t, req.FixtureDir, "node_modules/file", file4K)
	setPnpmStoreEnv(t, storeFiles)
	return nil
}
```