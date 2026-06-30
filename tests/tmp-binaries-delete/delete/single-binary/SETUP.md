# Scenario

**Leaf**: deleting one scanned binary removes the file and reports matching `freedSize`

## Steps

1. Create `~/Projects/del-one/.git` with Mach-O `bin/app` (120 bytes payload).
2. Scan then delete the scanned path.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/del-one")
	writeMachO(t, app, "bin/app")
	req.ScanFirst = true
	req.Op = "delete-binaries"
	return nil
}
```
