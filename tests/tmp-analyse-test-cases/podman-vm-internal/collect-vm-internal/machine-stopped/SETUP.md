# Scenario

**Feature**: stopped Podman machine skips vmInternal collection

```
# machine list shows stopped -> CollectPodmanVmInternal returns nil
podman machine list -> stopped -> skip SSH -> nil VmInternal
```

## Preconditions

- Machine list fixture reports no running machine.

## Steps

1. Set `req.Op` to `collect-vm-machine-stopped`.
2. Load stopped machine list fixture.

## Context

- Frontend hides vm-internal-section when machine not running.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "collect-vm-machine-stopped"
	req.MachineListJSON = loadMachineListFixture(t, "../testdata/machine-stopped.json")
	req.MachineRunning = false
	return nil
}
```