# Scenario

**Decision**: successful binary deletion

```
delete -> single | multiple
```

## Preconditions

- Paths come from the current scan session.
- Files are regular binary files on disk.

## Steps

1. Create fixture binaries in git repos.
2. Scan then delete selected paths.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "delete-binaries"
	return nil
}
```
