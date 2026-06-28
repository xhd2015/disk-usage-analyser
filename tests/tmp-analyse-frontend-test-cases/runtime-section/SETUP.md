# Scenario

**Feature**: Docker/Podman runtime section on cards

```
# location event carries runtimeItems; frontend renders runtime-section rows
CollectRuntimeStats -> location event -> runtime-section / runtime-row-{idx}
```

## Preconditions

- Docker and Podman cards may show a Runtime section below filesystem breakdown.
- Section hidden when runtimeItems empty or CLI unavailable.

## Steps

1. Descendant leaf runs playwright script after scan completes.

## Context

- Gracefully absent when runtime stats unavailable.

```go
func Setup(t *testing.T, req *Request) error {
	_ = req.ScriptFile
	return nil
}
```