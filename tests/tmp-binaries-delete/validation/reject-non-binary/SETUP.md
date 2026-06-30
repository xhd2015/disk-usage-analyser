# Scenario

**Leaf**: scanned binary path that no longer classifies as binary is rejected

## Steps

1. Create `~/Projects/bad-type/.git` with Mach-O `bin/good`.
2. Scan to record `bin/good` in the session.
3. Overwrite `bin/good` with plain text before delete.
4. Attempt delete of the scanned path.

## Context

- Re-validation uses the same `classifyFile` logic as scan.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/bad-type")
	writeMachO(t, app, "bin/good")
	req.ScanFirst = true
	req.OverwriteBeforeDelete = map[string]string{
		"Projects/bad-type/bin/good": "plain text not a binary\n",
	}
	req.Op = "delete-binaries"
	return nil
}
```