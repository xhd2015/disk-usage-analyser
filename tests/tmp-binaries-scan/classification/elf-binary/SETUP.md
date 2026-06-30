# Scenario

**Leaf**: ELF fixture is classified as `kind=elf`

## Steps

1. Create `~/Projects/elf-app/.git`.
2. Write ELF bytes to `bin/app`.
3. Run `binaries-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/elf-app")
	writeELF(t, app, "bin/app")
	req.Op = "binaries-scan"
	return nil
}
```
