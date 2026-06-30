# Scenario

**Leaf**: source and text files inside a git repository are not reported

## Preconditions

- The repository contains only `.go` source and `.txt` files.

## Steps

1. Create `~/Projects/text-only/.git`.
2. Write `main.go` and `README.txt`.
3. Run `scan`.

## Context

- File extension alone never creates a hit; detection must identify a supported binary.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/text-only")
	writeText(t, app, "main.go", "package "+"main\nfunc "+"main() {}\n")
	writeText(t, app, "README.txt", "not a binary\n")
	req.Args = []string{"scan"}
	return nil
}
```
