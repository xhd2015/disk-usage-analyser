# Scenario

**Leaf**: unreadable subtree is skipped gracefully while partial results stream

## Preconditions

- The fixture contains one readable repository and one unreadable subtree.

## Steps

1. Create `~/Projects/ok/.git` with a Mach-O binary.
2. Create `~/Projects/locked/.git` with a Mach-O binary.
3. Remove permissions from `locked`.
4. Run `scan`.

## Context

- The test restores permissions with cleanup so temporary directories can be removed.

```go
import "os"

func Setup(t *testing.T, req *Request) error {
	ok := repo(t, req.HomeDir, "Projects/ok")
	writeMachO(t, ok, "bin/ok")
	locked := repo(t, req.HomeDir, "Projects/locked")
	writeMachO(t, locked, "bin/locked")
	if err := os.Chmod(locked, 0000); err != nil {
		return err
	}
	t.Cleanup(func() {
		_ = os.Chmod(locked, 0755)
	})
	req.Args = []string{"scan"}
	return nil
}
```
