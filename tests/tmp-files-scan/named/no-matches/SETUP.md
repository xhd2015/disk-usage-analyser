# Scenario

**Leaf**: no directory or file matches `--name`; summary shows zero named entries

## Preconditions

- A git repository exists but contains no file or directory matching the `--name` value.

## Steps

1. Create `~/Projects/app/.git`.
2. Create `app/src/main.go` (a plain text file, not a match for --name).
3. Run `scan --name nonexistent`.

## Context

- When no entries match `--name`, zero named hits are reported, but the summary still includes the named count.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/app")
	writeText(t, app, "src/main.go", "package "+"main\n")
	req.Args = []string{"scan", "--name", "nonexistent"}
	return nil
}
```
