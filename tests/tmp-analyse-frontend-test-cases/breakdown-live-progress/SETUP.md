# Scenario

**Feature**: live breakdown progress during multi-path scan

```
# progress SSE updates breakdown row sizes before location event completes
Start Scan -> progress events -> breakdown-size-0/1 update -> done-badge
```

## Preconditions

- Go card has at least two breakdown rows (pkg/mod + go-build cache).
- Progress events include breakdownIndex and breakdownSize fields.

## Steps

1. Descendant leaf runs playwright script against Go card during scan.

## Context

- Verifies breakdown rows update before the done-badge appears.

```go
func Setup(t *testing.T, req *Request) error {
	if req.ScriptFile == "" {
		req.ScriptFile = "go-multi-path/breakdown-live-progress.js"
	}
	return nil
}
```