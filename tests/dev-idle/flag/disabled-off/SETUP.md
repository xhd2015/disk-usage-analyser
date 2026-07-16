# Scenario

**Leaf**: `--dev-idle-life off` disables idle life (zero duration; `0` is equivalent)

```
run.RunWithOptions(["--dev", "--dev-idle-life", "off"])
  -> StartServer(ServerOptions{Dev: true, DevIdleLife: 0})
```

## Steps

1. Pass `["--dev", "--dev-idle-life", "off"]` to `run.RunWithOptions`.
2. Capture `ServerOptions` via fake `StartServer`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Scenario = "flag/disabled-off"
	req.Args = []string{"--dev", "--dev-idle-life", "off"}
	return nil
}
```