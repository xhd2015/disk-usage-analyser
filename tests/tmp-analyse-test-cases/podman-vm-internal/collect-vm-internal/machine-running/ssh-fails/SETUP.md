# Scenario

**Feature**: SSH failure omits vmInternal without scan error

```
# running machine but SSH fails -> nil VmInternal, no collection error
CollectPodmanVmInternal -> SSH error -> nil VmInternal gracefully
```

## Preconditions

- Machine reports running but SSH commands fail.

## Steps

1. Set `req.Op` to `collect-vm-ssh-fails`.
2. Enable `MockSSHFail`.

## Context

- Overall scan must never fail due to VM stats unavailability.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "collect-vm-ssh-fails"
	req.MockSSHFail = true
	return nil
}
```