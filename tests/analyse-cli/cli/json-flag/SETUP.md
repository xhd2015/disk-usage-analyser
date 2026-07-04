# Scenario

**Leaf**: `--json` emits one JSON object with root, rows, and summary

## Steps

1. Create `sub/item.txt` (2048 B) in fixture.
2. Run `RunCLI --json <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	mkdir(t, req.FixtureDir, "sub")
	writeSizedFile(t, req.FixtureDir, "sub/item.txt", 2048)
	req.Mode = "cli"
	req.Args = []string{"--json", req.FixtureDir}
	return nil
}
```