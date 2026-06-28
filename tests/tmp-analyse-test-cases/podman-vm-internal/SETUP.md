# Scenario

**Feature**: Podman VM internal storage grouping

```
# Podman machine SSH -> du -sb / system df -> vmInternal + runtimeItems on location event
CollectPodmanVmInternal -> TmpVmInternal -> location event -> frontend vm-internal-section
```

## Preconditions

- Grouping node for Podman VM internal storage leaves.
- Injectable `SetPodmanMachineRunner` for SSH and machine list mocks.

## Steps

1. Descendant leaf sets `req.Op` and scenario-specific fields.

## Context

- macOS-only collection; graceful omission on other platforms.
- See leaf ASSERT.md for concrete expectations.

```go
func Setup(t *testing.T, req *Request) error {
	if req.MockSSHOutputs == nil {
		req.MockSSHOutputs = make(map[string]string)
	}
	return nil
}
```