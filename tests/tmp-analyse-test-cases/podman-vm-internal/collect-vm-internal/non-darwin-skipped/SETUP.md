# Scenario

**Feature**: non-darwin platform skips VM internal collection

```
# GOOS != darwin -> CollectPodmanVmInternal returns nil without SSH
runtime.GOOS linux -> skip VM collection -> nil VmInternal
```

## Preconditions

- Platform override set to non-darwin (e.g. `linux`).

## Steps

1. Set `req.Op` to `collect-vm-non-darwin`.
2. Set `ForceGOOS` to `linux`.

## Context

- macOS-only feature; Linux/Windows host scans omit VM stats.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "collect-vm-non-darwin"
	req.ForceGOOS = "linux"
	req.MachineListJSON = loadMachineListFixture(t, "../testdata/machine-running.json")
	return nil
}
```