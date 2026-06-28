# Scenario

**Feature**: parse docker system df NDJSON

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- Parse docker NDJSON; filter keeps Images (12 total, 8 active, 8.3GB/1.5GB reclaimable) and Build Cache (34 total).

## Steps
1. Set req.Op and scenario fields.
2. Run executes the targeted server function.

## Context
- Platform-independent pure function or mock-runner test.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "parse-system-df-docker"
	req.FixtureFile = "testdata/docker-system-df.ndjson"
	return nil
}
```
