# Scenario

**Feature**: one 4K file cloned twice via APFS `cp -c`

## Steps

1. Create `original.bin` (4096 B).
2. Clone with `cp -c` to `clone1.bin` and `clone2.bin` (darwin only).
3. Run `analyse.Analyse`.

## Context

- Skipped on non-darwin via `runtime.GOOS` check.

```go
import (
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if runtime.GOOS != "darwin" {
		t.Skip("APFS clone fixtures require darwin")
	}
	original := writeSizedFile(t, req.FixtureDir, "original.bin", file4K)
	clone1 := filepath.Join(req.FixtureDir, "clone1.bin")
	clone2 := filepath.Join(req.FixtureDir, "clone2.bin")
	if err := exec.Command("cp", "-c", original, clone1).Run(); err != nil {
		t.Fatalf("cp -c original -> clone1: %v", err)
	}
	if err := exec.Command("cp", "-c", original, clone2).Run(); err != nil {
		t.Fatalf("cp -c original -> clone2: %v", err)
	}
	return nil
}
```