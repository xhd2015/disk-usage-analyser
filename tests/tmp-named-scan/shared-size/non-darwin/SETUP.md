# Scenario

**Leaf**: non-darwin platforms report `sharedSize=0` on named events

## Steps

1. Skip on darwin.
2. Create git repo with `node_modules/` and store/cache env pointing at fixture files.
3. Run `named-scan`.

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
	tempRoot := t.TempDir()
	app := nodeModulesRepo(t, req.HomeDir, "Projects/shared-app")
	writeSizedFile(t, app, "node_modules/pkg/file", file4K)

	storeFiles := pnpmStoreFilesDir(tempRoot)
	writePnpmStoreFile(t, storeFiles, file4K)
	cacheRoot := bunCacheRoot(tempRoot)
	writeBunCacheFile(t, cacheRoot, file4K)

	t.Setenv("DISK_USAGE_ANALYSER_PNPM_STORE", storeFiles)
	t.Setenv("DISK_USAGE_ANALYSER_BUN_CACHE", cacheRoot)
	return nil
}
```