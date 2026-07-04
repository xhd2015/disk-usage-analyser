# Scenario

**Leaf**: clone fixture with no env override and no default store under HOME → `pnpm_shared=0`

## Steps

1. Create fixture store file and clone into `node_modules/pkg/file` via `cp -c` (darwin).
2. Do **not** set `DISK_USAGE_ANALYSER_PNPM_STORE`.
3. Point `HOME` at an empty directory (no `Library/pnpm/store`).
4. Run `analyse.Analyse` on fixture root.

```go
import (
	"path/filepath"
	"runtime"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if runtime.GOOS != "darwin" {
		t.Skip("missing-store-env clone fixture requires darwin")
	}
	tempRoot := filepath.Dir(req.FixtureDir)
	storeFiles := pnpmStoreFilesDir(tempRoot)
	storeFile := writePnpmStoreFile(t, storeFiles, file4K)
	cloneStoreFileTo(t, storeFile, "node_modules/pkg/file", req.FixtureDir)
	t.Setenv("DISK_USAGE_ANALYSER_PNPM_STORE", "")
	isolateHomeWithoutPnpmStore(t)
	return nil
}
```