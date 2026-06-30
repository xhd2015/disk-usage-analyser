# Scenario

**Leaf**: explicit repeatable `--root` limits scanning to selected fixture subtree

## Preconditions

- The fake home contains one selected subtree and one unselected subtree.

## Steps

1. Create `~/selected/app/.git` with a Mach-O binary.
2. Create `~/unselected/app/.git` with a Mach-O binary.
3. Run `scan --root ~/selected`.

## Context

- Explicit roots override the default home root.

```go
func Setup(t *testing.T, req *Request) error {
	selected := mkdir(t, req.HomeDir, "selected")
	unselected := mkdir(t, req.HomeDir, "unselected")
	app := repo(t, selected, "app")
	writeMachO(t, app, "bin/selected-app")
	other := repo(t, unselected, "app")
	writeMachO(t, other, "bin/unselected-app")
	req.Args = []string{"scan", "--root", selected}
	return nil
}
```
