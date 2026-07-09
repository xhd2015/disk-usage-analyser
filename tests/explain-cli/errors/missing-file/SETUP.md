# Scenario

**Leaf**: non-existent PATH exits non-zero

## Steps

1. Pick a path under `t.TempDir()` that is not created.
2. Run `explain.RunCLI <missing-path>`.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	missing := filepath.Join(t.TempDir(), "does-not-exist")
	req.TargetPath = missing
	req.Args = []string{missing}
	return nil
}
```
