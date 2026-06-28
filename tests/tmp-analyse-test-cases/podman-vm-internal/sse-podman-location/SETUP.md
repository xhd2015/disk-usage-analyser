# Scenario

**Feature**: Podman location SSE event carries vmInternal and runtimeItems

```
# HandleTmpAnalyse scan -> Podman location event with vmInternal + runtimeItems
CollectPodmanVmInternal + CollectPodmanRuntimeViaSSH -> location event -> frontend
```

## Preconditions

- Grouping node for SSE integration of VM internal stats.

## Steps

1. Descendant leaf configures full mock chain for Podman location event.

## Context

- End-to-end handler test with injectable runners.

```go
func Setup(t *testing.T, req *Request) error {
	req.ForceGOOS = "darwin"
	return nil
}
```