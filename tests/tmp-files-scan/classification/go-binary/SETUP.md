# Scenario

**Leaf**: Go buildinfo takes precedence and reports `Kind=go`

## Preconditions

- The repository contains a real Go-built executable.

## Steps

1. Create `~/Projects/go-app/.git`.
2. Build a tiny Go main package to `bin/go-app`.
3. Run `scan`.

## Context

- On Darwin the Go binary may also be Mach-O; on Linux it may also be ELF. Build info decides `Kind=go`.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/go-app")
	writeGoBinary(t, app, "bin/go-app")
	req.Args = []string{"scan"}
	return nil
}
```
