# Scenario

**Leaf**: RunCLI with no PATH scans the current working directory

## Steps

1. Write `only.txt` (256 bytes) in the fixture root.
2. `Chdir` into the fixture directory.
3. Run `RunCLI` with empty args.

```go
import (
	"os"
)

func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "only.txt", 2<<20)
	if err := os.Chdir(req.FixtureDir); err != nil {
		t.Fatalf("chdir fixture: %v", err)
	}
	req.Args = nil
	return nil
}
```