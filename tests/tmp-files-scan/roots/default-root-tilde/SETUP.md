# Scenario

**Leaf**: omitted `--root` scans the fixture home and renders tilde paths

## Preconditions

- The CLI receives an injected fake home directory.
- The command omits `--root`.

## Steps

1. Create `~/Projects/default-app/.git`.
2. Add one Mach-O binary.
3. Run `scan` without root flags.

## Context

- This proves default root resolution uses `HomeDir`, not the real process home.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/default-app")
	writeMachO(t, app, "bin/default-app")
	req.Args = []string{"scan"}
	return nil
}
```
