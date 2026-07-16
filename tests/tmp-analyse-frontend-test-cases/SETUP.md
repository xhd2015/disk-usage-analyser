# Scenario

**Feature**: tmp-analyse frontend shared playwright harness

```
# React page loads cards, user starts scan, playwright-debug verifies DOM
User -> /tmp-analyse -> Start Scan -> SSE progress -> card/breakdown updates
```

## Preconditions

- The test harness starts its own Go server; no external server needed.
- `playwright-debug` CLI tool is installed and on PATH.

## Steps

1. Root Setup verifies playwright-debug is available.
2. Each leaf sets `req.ScriptFile` to its playwright script.

## Context

- Scripts use `data-testid` attributes and log `ELEM`, `PASS`, `FAIL` lines.
- Server startup may take up to 60 seconds (Vite dev server initialization).

```go
import (
	"os/exec"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if req.TeardownOnly {
		return nil
	}
	_, err := exec.LookPath("playwright-debug")
	if err != nil {
		return err
	}
	return nil
}
```