# Scenario

**Leaf**: non-existent directory argument exits 2

## Steps

1. Pick a path under `t.TempDir()` that is not created.
2. Run `RunCLI <missing-path>`.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	missing := filepath.Join(t.TempDir(), "does-not-exist")
	req.Mode = "cli"
	req.Args = []string{missing}
	return nil
}
```