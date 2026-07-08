# Scenario

**Feature**: mock simctl + on-disk sizes yields runtime items

```
# mock runner returns JSON; CollectSimulatorRuntimeStats walks mountPath sizes
SetSimulatorRuntimeCommandRunner(mock JSON) -> walk mountPath -> TmpRuntimeItem
```

## Preconditions

- Temp directory holds a file with known byte size at the mocked mountPath.
- Mock runner returns fixture JSON with mountPath pointing at temp dir.

## Steps

1. Create temp mount dir with 1024-byte file.
2. Substitute mountPath placeholder in fixture.
3. Set `req.Op` to `collect-simulator-mock-success`.

## Context

- Verifies end-to-end collection, not parse-only.

```go
import (
	"os"
	"strings"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "collect-simulator-mock-success"
	mountDir := writeSizedDir(t, t.TempDir(), "ios-18-5", 1024)
	data, err := os.ReadFile("testdata/simulator-runtimes.json")
	if err != nil {
		return err
	}
	req.MockOutput = strings.ReplaceAll(string(data), "__MOUNT_IOS_18_5__", mountDir)
	return nil
}
```