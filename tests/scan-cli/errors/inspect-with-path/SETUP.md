# Scenario

**Leaf**: `--inspect FILE` cannot be combined with positional PATH

## Steps

1. Write a minimal valid TreeResult JSON file.
2. Run `RunCLI --inspect <file> <positional-path>`.

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	scanRoot := filepath.Join(t.TempDir(), "scan-root")
	if err := os.MkdirAll(scanRoot, 0755); err != nil {
		t.Fatalf("mkdir scan-root: %v", err)
	}
	jsonPath := filepath.Join(t.TempDir(), "tree.json")
	tree, total := sampleInspectTree(scanRoot)
	writeTreeResultJSON(t, jsonPath, scanRoot, total, 0, 24, tree)
	req.Args = []string{"--inspect", jsonPath, req.FixtureDir}
	return nil
}
```
