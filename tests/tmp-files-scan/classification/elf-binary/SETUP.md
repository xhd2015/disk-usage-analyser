# Scenario

**Leaf**: ELF binary inside a git repository is reported as `elf`

## Preconditions

- The repository contains one ELF fixture and no Go build info.

## Steps

1. Create `~/Projects/elf-app/.git`.
2. Write an ELF executable stub under `bin/elf-app`.
3. Run `scan`.

## Context

- ELF classification uses `TypeDesc` prefix matching after buildinfo fails.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/elf-app")
	writeELF(t, app, "bin/elf-app")
	req.Args = []string{"scan"}
	return nil
}
```
