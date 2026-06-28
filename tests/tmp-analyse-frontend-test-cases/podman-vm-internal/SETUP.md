# Scenario

**Feature**: Podman card Inside VM section after scan

```
# location event vmInternal -> frontend vm-internal-section rows
CollectPodmanVmInternal -> location event -> vm-internal-section / vm-internal-row-{idx}
```

## Preconditions

- Podman card may show an "Inside VM" section below host filesystem info.
- Section hidden when machine not running or vmInternal unavailable.

## Steps

1. Descendant leaf runs playwright script after scan completes.

## Context

- Card header size remains host filesystem total (unchanged summary bar).

```go
func Setup(t *testing.T, req *Request) error {
	_ = req.ScriptFile
	return nil
}
```