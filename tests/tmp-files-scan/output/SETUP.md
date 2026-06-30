# Scenario

**Decision**: output mode and timing

```
output -> human | json | streaming | tilde paths
```

## Preconditions

- Human output is default.
- `--json` emits NDJSON objects for hits, then the same text summary line.
- Every hit is flushed before the final summary.

## Steps

1. Create repositories with deterministic binary hits.
2. Select human or JSON output.
3. Inspect stdout text and structured fields.

## Context

- The streaming leaf uses a recording writer to detect whether the first write is a hit line.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{"scan"}
	}
	return nil
}
```
