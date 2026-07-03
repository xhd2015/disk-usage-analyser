# Scenario

**Decision**: multiple `--name` values match directories

```
scan --name=node_modules --name=vendor -> dir match node_modules | dir match vendor -> NamedHits
```

## Preconditions

- Two or more `--name` values are specified.
- Directories with those basenames exist under git repositories.

## Steps

1. Build fixtures with directories matching each name value.
2. Run scan with multiple `--name` flags.
3. Verify all named entries are reported, and non-named ignored dirs are still skipped.

## Context

- Ignores interact with names: only basenames in `--name` override `ignoredDirBasenames`; others remain skipped.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{"scan"}
	}
	return nil
}
```
