# Scenario

**Leaf**: omitting PATH, `--kind`, and `--all-kinds` exits non-zero

## Steps

1. Run `explain.RunCLI` with empty args (no PATH, no `--kind`, no `--all-kinds`, no other flags).

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{}
	return nil
}
```
