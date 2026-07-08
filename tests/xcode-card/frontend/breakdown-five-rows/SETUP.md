# Scenario

**Feature**: Xcode card shows five breakdown rows after locations load

```
# locations event carries five breakdown items for xcode
locations(xcode) -> breakdown-row-0..4 on card
```

## Preconditions

- Xcode card is detected on the test machine.

## Steps

1. Set `req.ScriptFile` to `breakdown-five-rows.js`.
2. Script verifies breakdown-row-0 through breakdown-row-4 exist.

## Context

- Row 4 label should include DocumentationCache path.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "breakdown-five-rows.js"
	return nil
}
```