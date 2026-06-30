# Scenario

**Decision**: delete path validation and safety

```
validation -> non-binary | directory | not in scan | already deleted
```

## Preconditions

- Server re-validates each path before unlink.
- Only paths from the current scan session are accepted.

## Steps

1. Build fixture isolating one rejection case.
2. Scan (unless testing no-scan rejection) then POST delete.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "delete-binaries"
	return nil
}
```
