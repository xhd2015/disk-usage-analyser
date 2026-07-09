# Scenario

**Leaf**: `--inspect` missing FILE exits non-zero

## Steps

1. Point `--inspect` at a path that does not exist under a temp dir.
2. Run `RunCLI --inspect <missing>`.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	missing := filepath.Join(t.TempDir(), "no-such-tree.json")
	req.Args = []string{"--inspect", missing}
	return nil
}
```
