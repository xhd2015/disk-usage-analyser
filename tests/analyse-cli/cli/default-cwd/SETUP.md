# Scenario

**Leaf**: RunCLI with no DIR analyses the current working directory

## Steps

1. Place `only.txt` (1024 B) in the fixture root.
2. `Chdir` into the fixture directory.
3. Run `RunCLI` with empty args.

```go
import (
	"os"
)

func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "only.txt", 1024)
	if err := os.Chdir(req.FixtureDir); err != nil {
		t.Fatalf("chdir fixture: %v", err)
	}
	req.Mode = "cli"
	req.Args = nil
	return nil
}
```