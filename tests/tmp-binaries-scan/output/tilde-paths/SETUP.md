# Scenario

**Leaf**: binary hit paths under home use `~/` prefix

## Steps

1. Create `~/Projects/tilde-bin/.git` with Mach-O binary.
2. Run `binaries-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/tilde-bin")
	writeMachO(t, app, "bin/app")
	req.Op = "binaries-scan"
	return nil
}
```
