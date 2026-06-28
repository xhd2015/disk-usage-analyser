# Scenario

**Feature**: tmp-analyse backend shared test harness

```
# handler discovers locations, scans paths, streams SSE progress and location events
Client -> HandleTmpAnalyse -> DiscoverLocations -> ScanWithProgress -> SSE events

# runtime stats collected via container CLI for Docker/Podman cards
HandleTmpAnalyse -> CollectRuntimeStats -> ParseSystemDFJSON -> location event runtimeItems
```

## Preconditions

- The server package provides TmpLocation, TmpBreakdownItem, TmpRuntimeItem, and TmpSummary types.
- DiscoverLocations, CalculateSize, ScanWithProgress, BuildSummary, and BuildProgressPayload exist.
- Mock filesystems (testing/fstest.MapFS) are used for deterministic file size tests.
- SSE format tests use httptest to capture event stream output.

## Steps

1. Verify server types compile and have expected fields.
2. Each leaf test sets `req.Op` to select the Run dispatch branch.

## Context

- Tests exercise the tmp-analyse feature of disk-usage-analyser.
- DiscoverLocations returns 5 core macOS locations plus 17+ software cache/log locations.
- Platform-specific full scans require darwin; parser and payload tests are platform-independent.

```go
import (
	"testing"

	"disk-usage-analyser/server"
)

func Setup(t *testing.T, req *Request) error {
	_ = server.TmpLocation{Path: "x", Label: "x", Category: "trash", Detected: true}
	_ = server.TmpSummary{}
	return nil
}
```