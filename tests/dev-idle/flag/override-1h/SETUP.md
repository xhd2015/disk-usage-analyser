# Scenario

**Leaf**: explicit `--dev-idle-life 1h` overrides the default

```
run.RunWithOptions(["--dev", "--dev-idle-life", "1h"])
  -> StartServer(ServerOptions{Dev: true, DevIdleLife: 1h})
```

## Steps

1. Pass `["--dev", "--dev-idle-life", "1h"]` to `run.RunWithOptions`.
2. Capture `ServerOptions` via fake `StartServer`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Scenario = "flag/override-1h"
	req.Args = []string{"--dev", "--dev-idle-life", "1h"}
	return nil
}
```