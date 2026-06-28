# Scenario

**Feature**: Go card breakdown rows update live during scan

```
# progress events patch breakdown-size-0 and breakdown-size-1 before done-badge
Start Scan -> progress -> breakdown-size updates -> done-badge
```

## Preconditions

- Go card is detected on the test machine.
- Go card has breakdown-items with at least two rows.

## Steps

1. Set req.ScriptFile to breakdown-live-progress.js.
2. Script starts scan and polls breakdown-size-0/1 while scanning-badge is visible.

## Context

- At least one breakdown row must show non-zero size before done-badge.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "breakdown-live-progress.js"
	return nil
}
```