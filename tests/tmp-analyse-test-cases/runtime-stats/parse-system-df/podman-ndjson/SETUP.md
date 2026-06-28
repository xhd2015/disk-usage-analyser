# Scenario

**Feature**: parse podman system df NDJSON with Image type alias

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- Podman Image row normalized to Images; Build Cache retained.

## Steps
1. Set req.Op and scenario fields.
2. Run executes the targeted server function.

## Context
- Platform-independent pure function or mock-runner test.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "parse-system-df-podman"
	req.FixtureFile = "testdata/podman-system-df.ndjson"
	return nil
}
```
