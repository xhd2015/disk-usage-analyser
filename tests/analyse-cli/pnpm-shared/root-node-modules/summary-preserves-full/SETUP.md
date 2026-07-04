# Scenario

**Bug**: summary `pnpm_shared` preserved when analyse root is `node_modules`

```
# fixture: store clone in node_modules/pkg/file; analyse node_modules/ directly
store files/ -> cp -c -> node_modules/pkg/file -> analyse(node_modules) -> summary pnpm_shared=4096
```

## Steps

1. Create fixture store `files/aa/<hash>` (4096 B).
2. Clone store file into `node_modules/pkg/file` via `cp -c`.
3. Set `DISK_USAGE_ANALYSER_PNPM_STORE` to fixture `files/` dir.
4. Point analyse root at `node_modules` (not parent fixture).
5. Run `analyse.Analyse` on the `node_modules` directory.

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
	cloneStoreFileTo(t, storeFile, "node_modules/pkg/file", fixtureParent)
	setPnpmStoreEnv(t, storeFiles)
	req.FixtureDir = filepath.Join(fixtureParent, "node_modules")
	return nil
}
```