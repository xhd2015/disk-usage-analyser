# Scenario

**Leaf**: paths under home are rendered with `~` in human output and result display fields

## Preconditions

- The repo and binary both live under the injected fixture home.

## Steps

1. Create `~/Projects/tilde-app/.git`.
2. Add one Mach-O binary.
3. Run `scan`.

## Context

- The raw absolute fixture home path must not leak into user-facing output fields.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/tilde-app")
	writeMachO(t, app, "bin/tilde-app")
	req.Args = []string{"scan"}
	return nil
}
```
