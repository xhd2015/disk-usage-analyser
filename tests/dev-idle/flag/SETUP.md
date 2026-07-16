# Scenario

**Decision**: `--dev-idle-life` CLI flag plumbing into `run.ServerOptions`

```
# Parse global server flags before StartServer
run.RunWithOptions(argv) -> lessflags (--dev, --dev-idle-life) -> ServerOptions

# --dev-idle-life only applies when --dev is set
--dev + omitted flag -> DevIdleLife = 10m (default)
--dev + --dev-idle-life <duration> -> DevIdleLife = parsed duration
--dev + --dev-idle-life off|0 -> DevIdleLife = 0 (disabled)
no --dev + --dev-idle-life <duration> -> DevIdleLife = 0 (ignored)
--dev + invalid duration -> error, StartServer not called

# Help surface
disk-usage-analyser -h -> documents --dev-idle-life
```

## Preconditions

- Target package: `disk-usage-analyser/run` (`ServerOptions.DevIdleLife` not yet implemented).
- Leaves use `run.RunWithOptions` with a fake `StartServer` hook that captures `ServerOptions`.
- No real web server or `Serve()` watchdog wiring in P2.

## Steps

1. Set `req.Scenario` to `flag/<leaf-slug>`.
2. Set `req.Args` to the argv slice for `run.RunWithOptions`.
3. Assert captured `DevIdleLife`, errors, or help stdout per leaf.

```go
import "bytes"

func Setup(t *testing.T, req *Request) error {
	if req.Stdout == nil {
		req.Stdout = &bytes.Buffer{}
	}
	if req.Stderr == nil {
		req.Stderr = &bytes.Buffer{}
	}
	return nil
}
```