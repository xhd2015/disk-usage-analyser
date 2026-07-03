# Scenario

**Decision**: directory basename matches a `--name` value

```
# directory with matching basename: compute recursive size, report NamedHit, SkipDir
repo walk -> dir.Name() in names -> computeDirSize(dir, names) -> NamedHit -> SkipDir

# size computation excludes nested same-name dirs (separate hits)
computeDirSize -> walk dir -> sum file sizes, skip subdirs whose basename is in names
```

## Preconditions

- The --name value matches a directory basename under a git repository.
- Nested directories with matching basenames are excluded from the outer size.

## Steps

1. Create a git repository with named directories containing files.
2. Run scan with the matching `--name` flag.
3. Verify NamedHit count, size, and output.

## Context

- Named directory hits use `SkipDir` so the subtree is not walked for binary detection.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{"scan"}
	}
	return nil
}
```
