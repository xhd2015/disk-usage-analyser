# Scenario

**Feature**: CollectSimulatorRuntimeStats via injectable simctl runner

```
# CollectSimulatorRuntimeStats invokes xcrun simctl runtime list -j
SetSimulatorRuntimeCommandRunner -> CollectSimulatorRuntimeStats -> []TmpRuntimeItem
```

## Preconditions

- Grouping node for mock-success and command-fails leaves.

## Steps

1. Descendant leaf sets `req.Op` and scenario fields.

## Context

- See leaf ASSERT.md for concrete expectations.

```go
func Setup(t *testing.T, req *Request) error {
	if req.FixtureFile == "" {
		req.FixtureFile = "testdata/simulator-runtimes.json"
	}
	return nil
}
```