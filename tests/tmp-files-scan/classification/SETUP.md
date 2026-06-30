# Scenario

**Decision**: file classification

```
classification -> go | macho | elf | non-binary
```

## Preconditions

- Go build info takes precedence over Mach-O or ELF container type.
- Mach-O and ELF are reported only when file type detection identifies binary magic.
- Text and source files are ignored.

## Steps

1. Create one repository with classification-specific files.
2. Run `scan --go-binaries` or default scan.

## Context

- Today `--go-binaries` is the only type selector and is equivalent to the default scan.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{"scan"}
	}
	return nil
}
```
