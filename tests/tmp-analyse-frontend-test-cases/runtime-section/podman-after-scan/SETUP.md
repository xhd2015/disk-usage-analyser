# Scenario

**Feature**: Podman card shows runtime section after scan when stats available

```
# Podman location event includes runtimeItems -> runtime-section rendered
CollectRuntimeStats(podman) -> location event -> runtime-section / runtime-row-0
```

## Preconditions

- Podman card is detected (~/.local/share/containers exists).
- When podman CLI returns stats, runtime-section appears with runtime-row-0.

## Steps

1. Set req.ScriptFile to podman-runtime-section.js.
2. Script completes scan and checks runtime-section or graceful absence.

## Context

- Test passes whether runtime stats are available or gracefully omitted.
- On macOS, Podman runtime stats may be collected inside the VM via SSH rather than host `podman system df`.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "podman-runtime-section.js"
	return nil
}
```