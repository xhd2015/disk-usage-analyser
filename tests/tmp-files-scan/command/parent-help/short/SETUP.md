# Scenario

**Bug**: `tmp-files -h` (parent level, short form) prints scan usage and exits 0

```
tmp-files -h -> usage (same as empty / scan -h / --help) -> exit 0
```

## Steps

1. Run `tmpfiles.RunCLI` with args `[]string{"-h"}` (no `scan` token).
2. Capture stdout, exit code, and error.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "parent-help"
	req.Args = []string{"-h"}
	return nil
}
```
