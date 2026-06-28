# Scenario

**Feature**: VM internal collection when Podman machine is running

```
# running machine -> SSH du on storage + overlay -> two TmpVmStorageItem rows
CollectPodmanVmInternal -> machine running -> SSH du -> items populated
```

## Preconditions

- Machine list reports running state.
- SSH `du -sb` mocked for storage paths.

## Steps

1. Descendant leaf sets success or SSH-failure mock.

## Context

- Primary happy path and SSH error path share running-machine precondition.

```go
func Setup(t *testing.T, req *Request) error {
	req.MachineListJSON = loadMachineListFixture(t, "../../testdata/machine-running.json")
	req.MachineRunning = true
	return nil
}
```