# Scenario

**Feature**: inner runtime collection fails gracefully

```
# SSH system df fails inside VM -> empty runtimeItems, no error
CollectPodmanRuntimeViaSSH -> SSH error -> nil/empty runtimeItems
```

## Preconditions

- Machine running but inner `podman system df` SSH fails.

## Steps

1. Set `req.Op` to `collect-runtime-ssh-fails`.
2. Enable `MockSSHFail` for runtime command path.

## Context

- Podman card still shows host filesystem scan; runtime section hidden.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "collect-runtime-ssh-fails"
	req.MockSSHFail = true
	return nil
}
```