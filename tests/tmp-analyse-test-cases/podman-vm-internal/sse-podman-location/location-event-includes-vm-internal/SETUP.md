# Scenario

**Feature**: Podman location SSE includes vmInternal and runtimeItems

```
# mocked VM collection + inner runtime -> Podman location event JSON
HandleTmpAnalyse -> Podman location event -> vmInternal + runtimeItems present
```

## Preconditions

- Mock machine running, SSH du fixtures, and inner system df fixture.
- Podman card detected in scan flow.

## Steps

1. Set `req.Op` to `sse-podman-location`.
2. Wire all SSH mock outputs.

## Context

- Validates JSON fields on location event for frontend rendering.

```go
import (
	"os"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "sse-podman-location"
	req.MachineListJSON = loadMachineListFixture(t, "../../collect-vm-internal/testdata/machine-running.json")
	req.MockSSHOutputs["/home/core/.local/share/containers/storage"] = loadDuFixture(t, "../../collect-vm-internal/testdata/storage-du.txt")
	req.MockSSHOutputs["/home/core/.local/share/containers/storage/overlay"] = loadDuFixture(t, "../../collect-vm-internal/testdata/overlay-du.txt")
	data, err := os.ReadFile("../../collect-runtime-via-ssh/mock-success/testdata/podman-system-df.ndjson")
	if err != nil {
		t.Fatalf("read runtime fixture: %v", err)
	}
	req.MockRuntimeOutput = string(data)
	return nil
}

func loadMachineListFixture(t *testing.T, name string) string {
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("read machine list fixture %s: %v", name, err)
	}
	return string(data)
}

func loadDuFixture(t *testing.T, name string) string {
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("read du fixture %s: %v", name, err)
	}
	return string(data)
}
```