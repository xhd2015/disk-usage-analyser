# Scenario

**Leaf**: deleting a scanned path that was removed before POST returns not-found failure

## Steps

1. Create `~/Projects/missing/.git` with Mach-O `bin/vanish`.
2. Scan to populate session with `bin/vanish`.
3. Remove `bin/vanish` from disk after scan.
4. Attempt delete of the scanned path.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/missing")
	writeMachO(t, app, "bin/vanish")
	req.ScanFirst = true
	req.RemoveAfterScan = []string{"~/Projects/missing/bin/vanish"}
	req.Op = "delete-binaries"
	return nil
}
```