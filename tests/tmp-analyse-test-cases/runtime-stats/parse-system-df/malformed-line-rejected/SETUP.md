# Scenario

**Feature**: reject malformed NDJSON lines

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- Malformed NDJSON line returns parse error.

## Steps
1. Set req.Op and scenario fields.
2. Run executes the targeted server function.

## Context
- Platform-independent pure function or mock-runner test.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "parse-system-df-malformed"
	req.FixtureFile = "testdata/bad.ndjson"
	return nil
}
```
