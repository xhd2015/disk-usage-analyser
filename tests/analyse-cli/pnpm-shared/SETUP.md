# Scenario

**Feature**: pnpm store clone matching for `pnpm_shared` column

```
analyse -> walk node_modules -> match APFS clone keys against pnpm store index -> unique bytes per key -> bubble up to parent rows
```

```

## Preconditions

- Darwin store-match leaves use `cp -c` (APFS clone) like `clones/apfs-three-refs`.
- Tests point the store at an isolated fixture via `DISK_USAGE_ANALYSER_PNPM_STORE` (absolute path to `files/`).
- Each store and `node_modules` file uses 4096-byte payloads for deterministic 4K metrics.

## Context

- Clone group key: `doc_id` when non-zero, else `clone_id`.
- Per `node_modules` walk: sum `file_size` once per unique matching clone key (Option A).
- Immediate-child row `pnpm_shared` = sum of all nested `node_modules` totals in that subtree.
- Missing or unreadable store → `pnpm_shared = 0` for all rows (no error).
- Non-darwin: `pnpm_shared` always `0`.

```go
import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

const pnpmHash = "abc123deadbeef"

func pnpmStoreFilesDir(tempRoot string) string {
	return filepath.Join(tempRoot, "store", "files")
}

func writePnpmStoreFile(t *testing.T, storeFilesDir string, size int64) string {
	t.Helper()
	rel := filepath.Join("aa", pnpmHash)
	return writeSizedFile(t, storeFilesDir, rel, size)
}

func setPnpmStoreEnv(t *testing.T, storeFilesDir string) {
	t.Helper()
	t.Setenv("DISK_USAGE_ANALYSER_PNPM_STORE", storeFilesDir)
}

func cloneStoreFileTo(t *testing.T, src string, destRel string, fixtureDir string) string {
	t.Helper()
	if runtime.GOOS != "darwin" {
		t.Fatal("APFS clone fixtures require darwin")
	}
	dest := filepath.Join(fixtureDir, filepath.FromSlash(destRel))
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		t.Fatalf("mkdir clone parent %s: %v", destRel, err)
	}
	if err := exec.Command("cp", "-c", src, dest).Run(); err != nil {
		t.Fatalf("cp -c %s -> %s: %v", src, destRel, err)
	}
	return dest
}

func isolateHomeWithoutPnpmStore(t *testing.T) string {
	t.Helper()
	home := filepath.Join(t.TempDir(), "empty-home")
	if err := os.MkdirAll(home, 0755); err != nil {
		t.Fatalf("mkdir home: %v", err)
	}
	t.Setenv("HOME", home)
	return home
}

func Setup(t *testing.T, req *Request) error {
	req.Mode = "analyse"
	return nil
}
```