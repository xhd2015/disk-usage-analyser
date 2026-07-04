# Scenario

**Feature**: analyse root is `node_modules` — summary and breakdown both reflect full upfront scan

```
# root walk computes full bun_shared via upfront node_modules scan
analyse node_modules/ -> walkSubtree -> summary bun_shared=4096

# per-package child row carries cache-backed bytes; child sum equals summary
pkg row bun_shared=4096 -> breakdown reconciles with summary
```

## Preconditions

- Requires darwin (`cp -c`); skipped on other platforms.
- Analyse target is the `node_modules` directory itself, not its parent project.
- `DISK_USAGE_ANALYSER_BUN_CACHE` points at fixture cache root.

## Context

- Per-package immediate children carry `bun_shared` bytes from the root `node_modules` scan when analyse root is `node_modules` (`summary-preserves-full`, `breakdown-attributes-per-package`).
- Summary `bun_shared` preserves the root `walkSubtree` total and equals the sum of immediate-child rows.

```go
import (
	"runtime"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if runtime.GOOS != "darwin" {
		t.Skip("root node_modules bun_shared alignment requires darwin")
	}
	return nil
}
```