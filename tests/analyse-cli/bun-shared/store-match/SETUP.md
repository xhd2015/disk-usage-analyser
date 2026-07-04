# Scenario

**Feature**: `node_modules` files APFS-cloned from fixture Bun cache index

```
cache walk -> CloneGroupKey set -> node_modules walk -> match keys -> bun_shared=Σ(unique file sizes)
```

```

## Preconditions

- Requires darwin (`cp -c`); skipped on other platforms.
- `DISK_USAGE_ANALYSER_BUN_CACHE` points at fixture cache root.

## Context

- Cache index is built lazily on first `node_modules` encountered during `Analyse`.
- Matching clone keys contribute their file size once per `node_modules` walk.

```go
import (
	"runtime"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if runtime.GOOS != "darwin" {
		t.Skip("bun cache clone fixtures require darwin")
	}
	return nil
}
```