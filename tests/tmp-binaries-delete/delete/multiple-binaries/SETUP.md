# Scenario

**Leaf**: deleting multiple scanned binaries removes all files and sums `freedSize`

## Steps

1. Create `~/Projects/del-multi/.git` with Mach-O `bin/a` and ELF `bin/b`.
2. Scan then delete all scanned paths.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/del-multi")
	writeMachO(t, app, "bin/a")
	writeELF(t, app, "bin/b")
	req.ScanFirst = true
	req.Op = "delete-binaries"
	return nil
}
```
