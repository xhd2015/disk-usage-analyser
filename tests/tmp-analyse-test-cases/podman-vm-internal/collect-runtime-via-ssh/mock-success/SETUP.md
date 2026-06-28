# Scenario

**Feature**: mock inner `podman system df` returns filtered runtime items

```
# SSH podman system df inside VM -> ParseSystemDFJSON -> Images + Build Cache
CollectPodmanRuntimeViaSSH -> mock NDJSON -> 2 runtimeItems
```

## Preconditions

- Mock SSH returns podman-system-df.ndjson fixture.
- Filtered to Images and Build Cache only.

## Steps

1. Set `req.Op` to `collect-runtime-ssh-success`.
2. Load NDJSON fixture into `MockRuntimeOutput`.

## Context

- Same filtering rules as host CollectRuntimeStats.

```go
import (
	"os"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "collect-runtime-ssh-success"
	data, err := os.ReadFile("testdata/podman-system-df.ndjson")
	if err != nil {
		t.Fatalf("read runtime fixture: %v", err)
	}
	req.MockRuntimeOutput = string(data)
	return nil
}
```