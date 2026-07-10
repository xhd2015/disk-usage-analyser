# Scenario

**Bug**: `tmp-files --help` (parent level, long form) prints scan usage and exits 0

```
tmp-files --help -> usage (same as empty / scan -h) -> exit 0
```

## Steps

1. Run `tmpfiles.RunCLI` with args `[]string{"--help"}` (no `scan` token).
2. Capture stdout, exit code, and error.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "parent-help"
	req.Args = []string{"--help"}
	return nil
}
```
