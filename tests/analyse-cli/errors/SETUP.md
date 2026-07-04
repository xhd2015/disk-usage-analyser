# Scenario

**Feature**: invalid inputs and usage errors

```
RunCLI -> validate path -> exit 2 on missing/invalid DIR
```

## Preconditions

- Missing or invalid paths must not panic; exit code 2 is returned.

## Context

- Unexpected walk failures use exit code 1 (out of scope for this leaf).

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	return nil
}
```