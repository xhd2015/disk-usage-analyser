# Scenario

**Leaf**: `--dev` without `--dev-idle-life` defaults idle life to 10 minutes

```
run.RunWithOptions(["--dev"]) -> StartServer(ServerOptions{Dev: true, DevIdleLife: 10m})
```

## Steps

1. Pass `["--dev"]` to `run.RunWithOptions`.
2. Capture `ServerOptions` via fake `StartServer`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Scenario = "flag/default-10m-when-dev"
	req.Args = []string{"--dev"}
	return nil
}
```