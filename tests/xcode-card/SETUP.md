# Scenario

**Feature**: Xcode card backend discovery and simulator runtime stats

```
# DiscoverLocations expands Xcode ExtraPaths; CollectSimulatorRuntimeStats -> runtimeItems
DiscoverLocations(xcode) -> ExtraPaths + breakdown rows
CollectSimulatorRuntimeStats -> ParseSimulatorRuntimeJSON -> TmpRuntimeItem -> SSE location event
```

## Preconditions

- Tests run against the `disk-usage-analyser/server` package.
- Simulator runtime collection uses injectable `SetSimulatorRuntimeCommandRunner`.

## Steps

1. Descendant leaf sets `req.Op` and scenario-specific fields.

## Context

- See leaf ASSERT.md for concrete expectations.

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	if req.HomeDir == "" {
		req.HomeDir = "/Users/testuser"
	}
	return nil
}

func writeSizedDir(t *testing.T, root string, rel string, size int64) string {
	dir := filepath.Join(root, rel)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	file := filepath.Join(dir, "data.bin")
	if err := os.WriteFile(file, make([]byte, size), 0644); err != nil {
		t.Fatalf("write %s: %v", file, err)
	}
	return dir
}
```