# Scenario

**Leaf**: `--inspect` with invalid JSON exits non-zero

## Steps

1. Write a file that is not valid TreeResult JSON.
2. Run `RunCLI --inspect <file>`.

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	bad := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(bad, []byte("not-json{{{"), 0644); err != nil {
		t.Fatalf("write bad json: %v", err)
	}
	req.Args = []string{"--inspect", bad}
	return nil
}
```
