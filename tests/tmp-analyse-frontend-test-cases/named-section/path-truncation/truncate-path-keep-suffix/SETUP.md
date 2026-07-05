# Scenario

**Feature**: truncatePathKeepSuffix pure function for prefix-hide path display

```
# long path + maxVisibleChars -> truncatePathKeepSuffix -> display string with full suffix
TruncatedPath <- truncatePathKeepSuffix(path, maxVisibleChars) -> ellipsis prefix + suffix
```

## Preconditions

- Node.js is installed and on PATH.
- `disk-usage-analyser-react/src/pathDisplay.ts` exports `truncatePathKeepSuffix`.
- Each leaf provides a JSON fixture under `testdata/`.

## Steps

1. Root Setup verifies `node` and `npx` are available.
2. Leaf Setup sets `req.Op` and `req.FixtureFile`.

```go
import (
	"os/exec"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if _, err := exec.LookPath("node"); err != nil {
		return err
	}
	if _, err := exec.LookPath("npx"); err != nil {
		return err
	}
	if req.FixtureFile == "" {
		req.FixtureFile = "testdata/fixture.json"
	}
	return nil
}
```