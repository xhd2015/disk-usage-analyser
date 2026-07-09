# Scenario

**Leaf**: `--at` alone focuses tree; no match section

## Steps

1. Write sample inspect JSON.
2. Run `RunCLI --inspect <file> --at <scanRoot>/big`.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	p := prepareSampleInspect(t, req)
	at := filepath.Join(p.ScanRoot, "big")
	req.Args = []string{"--inspect", p.JSONPath, "--at", at}
	return nil
}
```
