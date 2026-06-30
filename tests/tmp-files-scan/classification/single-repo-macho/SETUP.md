# Scenario

**Leaf**: Mach-O binary inside a git repository is reported as `macho`

## Preconditions

- The repository contains one Mach-O fixture and no other binaries.

## Steps

1. Create `~/Projects/macho-app/.git`.
2. Write a Mach-O executable stub under `bin/macho-app`.
3. Run `scan --go-binaries`.

## Context

- `--go-binaries` is currently equivalent to the default binary scan.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/macho-app")
	writeMachO(t, app, "bin/macho-app")
	req.Args = []string{"scan", "--go-binaries"}
	return nil
}
```
