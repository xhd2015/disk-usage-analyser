# Scenario

**Leaf**: `-h` prints explain usage documenting PATH, `--kind`, `--all-kinds`, `--json`, and `--color`

## Steps

1. Run `explain.RunCLI` with `-h` (no PATH).

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	req.Args = []string{"-h"}
	return nil
}
```
