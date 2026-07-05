# Scenario

**Leaf**: `bun.lockb` in project root yields `packageManager=bun`

## Steps

1. Create git repo with `bun.lockb` and `node_modules/`.
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := nodeModulesRepo(t, req.HomeDir, "Projects/bun-app")
	writeFile(t, app, "bun.lockb", []byte{0x00, 0x62, 0x75, 0x6e}, 0644)
	return nil
}
```