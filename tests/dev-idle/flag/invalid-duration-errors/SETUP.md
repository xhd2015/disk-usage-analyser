# Scenario

**Leaf**: invalid `--dev-idle-life` value returns a parse error

```
run.RunWithOptions(["--dev", "--dev-idle-life", "not-a-duration"]) -> error
```

## Steps

1. Pass invalid duration with `--dev`.
2. Expect non-nil error; `StartServer` must not run.

```go
func Setup(t *testing.T, req *Request) error {
	req.Scenario = "flag/invalid-duration-errors"
	req.Args = []string{"--dev", "--dev-idle-life", "not-a-duration"}
	return nil
}
```