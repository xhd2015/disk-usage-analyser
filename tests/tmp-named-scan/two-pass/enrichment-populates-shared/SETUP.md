# Scenario

**Leaf**: darwin `named_enriched` carries non-zero shared metrics from analyse

## Steps

1. Skip on non-darwin.
2. Create git repo with `node_modules/`.
3. Write pnpm store and bun cache fixture files (4096 B each).
4. `cp -c` store file into `node_modules/pnpm-pkg/file`.
5. Set `DISK_USAGE_ANALYSER_PNPM_STORE` and `DISK_USAGE_ANALYSER_BUN_CACHE`.
6. Run `named-scan`.

```go
import (
	"path/filepath"
	"runtime"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if runtime.GOOS != "darwin" {
		t.Skip("APFS clone fixtures require darwin")
	}
	tempRoot := t.TempDir()
	app := nodeModulesRepo(t, req.HomeDir, "Projects/enrich-app")

	storeFiles := pnpmStoreFilesDir(tempRoot)
	storeFile := writePnpmStoreFile(t, storeFiles, file4K)
	cpClone(t, storeFile, filepath.Join(app, "node_modules", "pnpm-pkg", "file"))

	cacheRoot := bunCacheRoot(tempRoot)
	writeBunCacheFile(t, cacheRoot, file4K)

	t.Setenv("DISK_USAGE_ANALYSER_PNPM_STORE", storeFiles)
	t.Setenv("DISK_USAGE_ANALYSER_BUN_CACHE", cacheRoot)
	return nil
}
```