# Scenario

**Leaf**: `--dev-idle-life` without `--dev` is ignored (DevIdleLife stays zero)

```
run.RunWithOptions(["--dev-idle-life", "1h"])
  -> StartServer(ServerOptions{Dev: false, DevIdleLife: 0})
```

## Steps

1. Pass `["--dev-idle-life", "1h"]` without `--dev`.
2. Capture `ServerOptions` via fake `StartServer`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Scenario = "flag/ignored-without-dev"
	req.Args = []string{"--dev-idle-life", "1h"}
	return nil
}
```