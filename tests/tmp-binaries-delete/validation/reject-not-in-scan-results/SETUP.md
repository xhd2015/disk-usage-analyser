# Scenario

**Leaf**: binary file that was not in the current scan session cannot be deleted

## Steps

1. Create `~/Projects/scope/.git` with Mach-O `bin/scanned` only.
2. Run scan (session contains `bin/scanned` only).
3. Add Mach-O `bin/secret` after scan completes.
4. Attempt delete of `~/Projects/scope/bin/secret`.

## Context

- Confirmed default: delete only paths from current scan results.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/scope")
	writeMachO(t, app, "bin/scanned")
	req.ScanFirst = false
	req.ExtraPaths = []string{"~/Projects/scope/bin/secret"}
	req.AddAfterScan = map[string][]byte{
		"Projects/scope/bin/secret": append([]byte{0xcf, 0xfa, 0xed, 0xfe, 0x0c, 0x00, 0x00, 0x01}, make([]byte, 96)...),
	}
	req.Op = "delete-binaries"
	return nil
}
```