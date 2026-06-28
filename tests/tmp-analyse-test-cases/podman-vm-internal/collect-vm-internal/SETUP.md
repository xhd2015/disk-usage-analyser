# Scenario

**Feature**: collect TmpVmInternal via Podman machine SSH

```
# machine list check -> SSH du -sb on storage paths -> TmpVmInternal items
CollectPodmanVmInternal -> podman machine list -> SSH du -> TmpVmInternal
```

## Preconditions

- Grouping node for VM internal collection outcomes.

## Steps

1. Descendant leaf configures machine list and SSH mock responses.

## Context

- Never fails overall scan; nil vmInternal when unavailable.

```go
import (
	"os"
)

const (
	vmStoragePath = "/home/core/.local/share/containers/storage"
	vmOverlayPath = "/home/core/.local/share/containers/storage/overlay"
)

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

func Setup(t *testing.T, req *Request) error {
	req.ForceGOOS = "darwin"
	return nil
}
```