# Scenario

**Leaf**: Mach-O fixture is classified as `kind=macho`

## Steps

1. Create `~/Projects/macho-app/.git`.
2. Write Mach-O bytes to `bin/app`.
3. Run `binaries-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/macho-app")
	writeMachO(t, app, "bin/app")
	req.Op = "binaries-scan"
	return nil
}
```
