# Scenario

**Leaf**: root `-h` documents `--dev-idle-life` (no server start)

```
disk-usage-analyser -h -> stdout contains --dev-idle-life; StartServer not called
```

## Steps

1. Pass `["-h"]` to `run.RunWithOptions`.
2. Assert help stdout mentions `--dev-idle-life`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Scenario = "flag/help-lists-flag"
	req.Args = []string{"-h"}
	return nil
}
```