# Scenario

**Bug**: nested `node_modules` under `.pnpm/` must not re-run `pnpmSharedForNodeModules`

```
# outer node_modules walk matches store clone keys once per tree
analyse project root -> walk outer node_modules -> pnpmSharedForNodeModules (dedup map)

# nested node_modules inside .pnpm must not trigger a second accumulation
walk .../.pnpm/.../node_modules -> (must skip separate pnpm_shared scan)
```

## Preconditions

- Requires darwin (`cp -c`); skipped on other platforms.
- `DISK_USAGE_ANALYSER_PNPM_STORE` points at fixture `store/files/`.

## Context

- Real pnpm layouts nest many `node_modules` dirs under `.pnpm/pkg@ver/node_modules/`.
- Each nested scan with a fresh dedup map double-counts overlapping clone keys.
- Fix: run `pnpmSharedForNodeModules` only for the outermost `node_modules` in a subtree.

```go
import (
	"runtime"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if runtime.GOOS != "darwin" {
		t.Skip("pnpm nested node_modules clone fixtures require darwin")
	}
	return nil
}
```