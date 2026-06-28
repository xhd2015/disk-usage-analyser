# Scenario

**Feature**: Docker location SSE event includes runtimeItems

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- Docker location event carries runtimeItems when mock runner returns fixture.

## Steps
1. Set req.Op and scenario fields.
2. Run executes the targeted server function.

## Context
- Platform-independent pure function or mock-runner test.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "sse-runtime-docker"
	req.FixtureFile = "testdata/docker-system-df.ndjson"
	return nil
}
```
