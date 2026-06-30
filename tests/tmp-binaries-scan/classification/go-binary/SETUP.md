# Scenario

**Leaf**: real Go-built executable is classified as `kind=go`

## Steps

1. Create `~/Projects/go-app/.git`.
2. Build tiny Go binary to `bin/go-app`.
3. Run `binaries-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/go-app")
	writeGoBinary(t, app, "bin/go-app")
	req.Op = "binaries-scan"
	return nil
}
```
