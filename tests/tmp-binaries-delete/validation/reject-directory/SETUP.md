# Scenario

**Leaf**: directory path is rejected and not removed

## Steps

1. Create `~/Projects/bad-dir/.git` with Mach-O `bin/app`.
2. Scan then attempt delete of repo root directory path.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/bad-dir")
	writeMachO(t, app, "bin/app")
	req.ScanFirst = false
	req.ExtraPaths = []string{"~/Projects/bad-dir"}
	req.Op = "delete-binaries"
	return nil
}
```
