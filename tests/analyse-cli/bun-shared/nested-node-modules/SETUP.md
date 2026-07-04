# Scenario

**Bug**: nested `node_modules` under `.pnpm/` must not re-run `bunSharedForNodeModules`

```
# outer node_modules walk matches cache clone keys once per tree
analyse project root -> walk outer node_modules -> bunSharedForNodeModules (dedup map)

# nested node_modules inside .pnpm must not trigger a second accumulation
walk .../.pnpm/.../node_modules -> (must skip separate bun_shared scan)
```

## Preconditions

- Requires darwin (`cp -c`); skipped on other platforms.
- `DISK_USAGE_ANALYSER_BUN_CACHE` points at fixture cache root.

## Context

- Real Bun layouts may nest `node_modules` dirs under `.pnpm/pkg@ver/node_modules/`.
- Each nested scan with a fresh dedup map double-counts overlapping clone keys.
- Fix: run `bunSharedForNodeModules` only for the outermost `node_modules` in a subtree.

```go
import (
	"runtime"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if runtime.GOOS != "darwin" {
		t.Skip("bun nested node_modules clone fixtures require darwin")
	}
	return nil
}
```