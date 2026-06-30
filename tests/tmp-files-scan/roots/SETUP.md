# Scenario

**Decision**: scan root source

```
roots -> default HOME | explicit --root
```

## Preconditions

- Default scans use the injected fixture home directory.
- Explicit roots override the default home root.

## Steps

1. Prepare repositories under fake home.
2. Select omitted `--root` or explicit `--root`.

## Context

- Paths below the fixture home render with `~`.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{"scan"}
	}
	return nil
}
```
