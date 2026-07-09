# Scenario

**Leaf**: `run.RunWithOptions` routes `scan --inspect` without starting the web server

## Steps

1. Write a fixed TreeResult JSON file.
2. Call `run.RunWithOptions` with `scan --inspect <file>`.

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	scanRoot := filepath.Join(t.TempDir(), "scan-root")
	if err := os.MkdirAll(scanRoot, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	jsonPath := filepath.Join(t.TempDir(), "tree.json")
	tree, total := sampleInspectTree(scanRoot)
	writeTreeResultJSON(t, jsonPath, scanRoot, total, 0, 24, tree)
	req.Mode = "dispatch"
	req.Args = []string{"scan", "--inspect", jsonPath}
	// Stash paths for assert via FixtureDir / unused Args side channel:
	// store scan root in FixtureDir for PATH checks.
	req.FixtureDir = scanRoot
	return nil
}
```
