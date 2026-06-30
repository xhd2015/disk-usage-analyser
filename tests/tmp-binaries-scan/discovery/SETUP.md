# Scenario

**Decision**: binary discovery scope

```
discovery -> multi repos | ignored vendor dirs
```

## Preconditions

- Only binaries inside discovered git repos are reported.
- Ignored basenames include `vendor`, `node_modules`, `.venv`.

## Steps

1. Build fixture tree isolating one discovery rule.
2. Run `binaries-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "binaries-scan"
	return nil
}
```
