# Scenario

**Bug**: per-package breakdown rows show `pnpm_shared=0` when analyse root is `node_modules`

Same gap as Bun: summary `pnpm_shared` is correct but package breakdown rows show `0 B`.

```
# fixture: pkg-a has store clone; pkg-b is regular; analyse node_modules/
store -> cp -c -> node_modules/pkg-a/file
write regular -> node_modules/pkg-b/file
analyse(node_modules) -> summary pnpm_shared=4096; pkg-a row should carry 4096
```

## Steps

1. Create fixture pnpm store file (4096 B).
2. Clone into `node_modules/pkg-a/file` via `cp -c`.
3. Write a non-cloned regular file in `node_modules/pkg-b/file`.
4. Set `DISK_USAGE_ANALYSER_PNPM_STORE` to fixture `files/` dir.
5. Analyse the `node_modules` directory (not parent).

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	tempRoot := filepath.Dir(req.FixtureDir)
	fixtureParent := req.FixtureDir
	storeFiles := pnpmStoreFilesDir(tempRoot)
	storeFile := writePnpmStoreFile(t, storeFiles, file4K)
	cloneStoreFileTo(t, storeFile, "node_modules/pkg-a/file", fixtureParent)
	writeSizedFile(t, fixtureParent, "node_modules/pkg-b/file", file4K)
	setPnpmStoreEnv(t, storeFiles)
	req.FixtureDir = filepath.Join(fixtureParent, "node_modules")
	return nil
}
```