# Scenario

**Decision**: regular file basename matches a `--name` value

```
# file with matching basename: report NamedHit with direct size, continue walking
repo walk -> file.Name() in names -> NamedHit(file.Size()) -> continue walk
```

## Preconditions

- A non-directory file whose basename matches a `--name` value exists under a git repository.

## Steps

1. Create a git repository with a regular file whose name matches the `--name` value.
2. Run scan with the matching flag.
3. Verify the file is reported as a named hit (not a binary hit).

## Context

- File matches do not trigger SkipDir; the walk continues to siblings.
- The file is not classified as a binary even if it had binary magic (but for test fixtures we use plain text).

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{"scan"}
	}
	return nil
}
```
