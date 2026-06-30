# Scenario

**Decision**: binary scan output shape

```
output -> summary totals | tilde paths
```

## Preconditions

- Summary aggregates repos, binaries, and total size.
- Home-contained paths render with `~/` prefix.

## Steps

1. Build fixture with known binary count.
2. Run `binaries-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "binaries-scan"
	return nil
}
```
