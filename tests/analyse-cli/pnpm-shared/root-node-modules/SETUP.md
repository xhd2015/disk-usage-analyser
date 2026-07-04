# Scenario

**Feature**: analyse root is `node_modules` — summary and breakdown both reflect full upfront scan

```
# root walk computes full pnpm_shared via upfront node_modules scan
analyse node_modules/ -> walkSubtree -> summary pnpm_shared=4096

# per-package child row carries store-backed bytes; child sum equals summary
pkg row pnpm_shared=4096 -> breakdown reconciles with summary
```

## Preconditions

- Requires darwin (`cp -c`); skipped on other platforms.
- Analyse target is the `node_modules` directory itself, not its parent project.
- `DISK_USAGE_ANALYSER_PNPM_STORE` points at fixture `store/files/`.

## Context

- Per-package immediate children carry `pnpm_shared` bytes from the root `node_modules` scan when analyse root is `node_modules` (`summary-preserves-full`, `breakdown-attributes-per-package`).
- Summary `pnpm_shared` preserves the root `walkSubtree` total and equals the sum of immediate-child rows.

```go
import (
	"runtime"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if runtime.GOOS != "darwin" {
		t.Skip("root node_modules pnpm_shared alignment requires darwin")
	}
	return nil
}
```