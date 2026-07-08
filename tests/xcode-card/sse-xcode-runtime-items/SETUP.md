# Scenario

**Feature**: Xcode location SSE event includes runtimeItems

```
# HandleTmpAnalyse attaches simulator runtimeItems on xcode location event
HandleTmpAnalyse -> case xcode -> CollectSimulatorRuntimeStats -> location event
```

## Preconditions

- Mock simctl runner returns fixture JSON with at least one runtime.
- Xcode location appears in SSE stream (catalog always includes xcode).

## Steps

1. Set `req.Op` to `sse-xcode-runtime`.
2. Run full HandleTmpAnalyse SSE handler with mock runner installed.

## Context

- Mirrors docker-location-event pattern for xcode category.

```go
import (
	"os"
	"strings"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "sse-xcode-runtime"
	mountDir := writeSizedDir(t, t.TempDir(), "ios-18-5", 512)
	data, err := os.ReadFile("testdata/simulator-runtimes.json")
	if err != nil {
		return err
	}
	req.MockOutput = strings.ReplaceAll(string(data), "__MOUNT_IOS_18_5__", mountDir)
	return nil
}
```