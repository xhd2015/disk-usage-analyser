# Scenario

**Feature**: Podman card shows vm-internal-section when VM stats available

```
# Podman location event includes vmInternal -> vm-internal-section rendered
CollectPodmanVmInternal -> location event -> vm-internal-section / vm-internal-row-0
```

## Preconditions

- Podman card is detected (~/.local/share/containers exists).
- When Podman machine running and stats collected, vm-internal-section appears.

## Steps

1. Set `req.ScriptFile` to podman-vm-internal-section.js.
2. Script completes scan and checks vm-internal-section or graceful absence.

## Context

- Test passes whether VM internal stats are available or gracefully omitted.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "podman-vm-internal-section.js"
	return nil
}
```