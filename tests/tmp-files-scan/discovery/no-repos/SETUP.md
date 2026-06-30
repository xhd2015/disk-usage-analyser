# Scenario

**Leaf**: root with no git repositories reports zero binaries

## Preconditions

- The fake home has ordinary directories but no `.git` directories.

## Steps

1. Write a text file and an outside-repo ELF file.
2. Run `scan`.

## Context

- No repo discovery means no binary discovery.

```go
func Setup(t *testing.T, req *Request) error {
	writeText(t, req.HomeDir, "notes/readme.txt", "hello\n")
	writeELF(t, req.HomeDir, "bin/outside")
	req.Args = []string{"scan"}
	return nil
}
```
