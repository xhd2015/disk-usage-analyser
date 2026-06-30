# Scenario

**Decision**: partial batch delete outcomes

```
batch -> some deleted, some failed in one request
```

## Preconditions

- Partial success is allowed in a single POST body.

## Steps

1. Build fixture with one valid scanned binary and one invalid path.
2. Scan then delete both in one request.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "delete-binaries"
	return nil
}
```