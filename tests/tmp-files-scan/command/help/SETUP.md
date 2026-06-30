# Scenario

**Leaf**: `tmp-files scan -h` documents scan flags

## Preconditions

- The scan command supports `-h` and `--help`.

## Steps

1. Run `tmpfiles.RunCLI` with `scan -h`.
2. Capture stdout and exit code.

## Context

- Help is a successful command, not a runtime scan.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "help"
	req.Args = []string{"scan", "-h"}
	return nil
}
```
