# Scenario

**Leaf**: lone `*.qcow2` file (not under an AVD) detects generic-qcow2

## Steps

1. Write `disks/vm-disk.qcow2` (500 bytes) with no surrounding AVD signatures.
2. Run `explain.RunCLI` on that file.

```go
func Setup(t *testing.T, req *Request) error {
	path := writeSizedFile(t, req.FixtureDir, "disks/vm-disk.qcow2", 500)
	req.TargetPath = path
	req.Args = []string{path}
	return nil
}
```
