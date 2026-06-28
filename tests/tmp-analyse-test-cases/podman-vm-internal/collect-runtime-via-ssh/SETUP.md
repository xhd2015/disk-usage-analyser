# Scenario

**Feature**: collect Podman runtime stats inside VM via SSH

```
# podman machine ssh -- podman system df -> ParseSystemDFJSON -> FilterRuntimeItems
CollectPodmanRuntimeViaSSH -> inner system df -> runtimeItems (Images, Build Cache)
```

## Preconditions

- Grouping node for inner-VM runtime collection outcomes.
- Uses existing ParseSystemDFJSON + FilterRuntimeItems pipeline.

## Steps

1. Descendant leaf configures SSH mock for inner `podman system df`.

## Context

- Replaces host `podman system df` for Podman card on darwin.

```go
import (
	"os"
)

func Setup(t *testing.T, req *Request) error {
	req.ForceGOOS = "darwin"
	req.MachineListJSON = loadMachineListFixture(t, "../../collect-vm-internal/testdata/machine-running.json")
	return nil
}

func loadMachineListFixture(t *testing.T, name string) string {
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("read machine list fixture %s: %v", name, err)
	}
	return string(data)
}
```