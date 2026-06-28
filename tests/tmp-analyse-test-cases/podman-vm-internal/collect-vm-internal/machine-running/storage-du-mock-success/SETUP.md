# Scenario

**Feature**: mock SSH du returns two VM storage items

```
# running machine + mocked du output -> Container storage + Overlay layers items
CollectPodmanVmInternal -> SSH du storage + overlay -> TmpVmInternal with 2 items
```

## Preconditions

- Mock SSH returns storage and overlay `du` fixtures.
- `MachineRunning` true on result.

## Steps

1. Set `req.Op` to `collect-vm-storage-success`.
2. Wire SSH mock outputs for both paths.

## Context

- Labels: "Container storage" and "Overlay layers".

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "collect-vm-storage-success"
	req.MockSSHOutputs[vmStoragePath] = loadDuFixture(t, "../../testdata/storage-du.txt")
	req.MockSSHOutputs[vmOverlayPath] = loadDuFixture(t, "../../testdata/overlay-du.txt")
	return nil
}
```