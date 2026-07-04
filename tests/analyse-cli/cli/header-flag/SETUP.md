# Scenario

**Leaf**: `--header` remains valid; TSV always prints column names before data rows

## Steps

1. Create empty fixture (summary only).
2. Run `RunCLI --header <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	req.Args = []string{"--header", req.FixtureDir}
	return nil
}
```