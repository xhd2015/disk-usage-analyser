# Scenario

**Feature**: parse `du -sb` stdout from Podman machine SSH

```
# podman machine ssh -- du -sb <path> -> ParseDuSBOutput -> byte total
CollectPodmanVmInternal -> SSH du output -> ParseDuSBOutput
```

## Preconditions

- Grouping node for `du -sb` output format variants.

## Steps

1. Descendant leaf sets `req.Op` and `req.DuOutput` or fixture path.

## Context

- Pure parsing; no machine state or platform guards.

```go
func Setup(t *testing.T, req *Request) error {
	return nil
}
```