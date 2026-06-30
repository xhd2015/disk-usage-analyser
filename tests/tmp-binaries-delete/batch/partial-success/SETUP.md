# Scenario

**Leaf**: batch with one valid scanned binary and one directory path partially succeeds

## Steps

1. Create `~/Projects/partial/.git` with Mach-O `bin/ok`.
2. Scan then POST delete with scanned `bin/ok` path and directory `~/Projects/partial`.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/partial")
	writeMachO(t, app, "bin/ok")
	req.ScanFirst = true
	req.ExtraPaths = []string{"~/Projects/partial"}
	req.Op = "delete-binaries"
	return nil
}
```