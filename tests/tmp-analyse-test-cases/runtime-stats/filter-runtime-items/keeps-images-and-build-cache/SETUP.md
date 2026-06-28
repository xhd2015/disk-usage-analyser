# Scenario

**Feature**: filter keeps only Images and Build Cache

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- Containers and Local Volumes excluded; Images and Build Cache kept.

## Steps
1. Set req.Op and scenario fields.
2. Run executes the targeted server function.

## Context
- Platform-independent pure function or mock-runner test.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "filter-runtime-items"
	req.FixtureFile = "testdata/docker-system-df.ndjson"
	return nil
}
```
