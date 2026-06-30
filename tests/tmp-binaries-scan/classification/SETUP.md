# Scenario

**Decision**: binary file classification

```
classification -> go | macho | elf
```

## Preconditions

- Go buildinfo takes precedence over container format (Mach-O/ELF).
- Text/source files are not reported.

## Steps

1. Place one classified binary per leaf inside a git repo.
2. Run `binaries-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "binaries-scan"
	return nil
}
```
