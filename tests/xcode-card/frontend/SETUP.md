# Scenario

**Feature**: Xcode card frontend UI verification

```
# React page renders xcode card breakdown, runtime section, cleanup popover
User -> /tmp-analyse -> playwright-debug checks data-testid elements
```

## Preconditions

- The test harness starts its own Go server; no external server needed.
- `playwright-debug` CLI tool is installed and on PATH.

## Steps

1. Root Setup verifies playwright-debug is available.
2. Each leaf sets `req.ScriptFile` to its playwright script.

## Context

- Leaves use `ui-automation` label; runtime-section skips when simctl unavailable.

```go
import (
	"os/exec"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, err := exec.LookPath("playwright-debug")
	if err != nil {
		return err
	}
	return nil
}
```