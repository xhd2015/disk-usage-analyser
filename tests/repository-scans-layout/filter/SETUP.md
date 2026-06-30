# Scenario

**Feature**: filter repository scan rows by size threshold

```
# unchecked <1M / <10M checkboxes hide items under threshold
repositoryScansLayout.filter* -> hide leaves and empty repos per rules
```

## Preconditions

- `ONE_MIB = 1048576` for binaries; `TEN_MIB = 10485760` for worktrees.
- Binary repo hidden when sum of visible binaries < 1 MiB.
- Worktree repo hidden when main checkout size < 10 MiB (linked filtered independently).

## Steps

1. Leaf sets `req.Op` and fixture `showUnder*` flag.

## Context

- See leaf ASSERT.md for visible/hidden paths.

```go
func Setup(t *testing.T, req *Request) error {
	if req.FixtureFile == "" {
		req.FixtureFile = "testdata/fixture.json"
	}
	return nil
}
```