# Scenario

**Feature**: unparseable `du -sb` output returns error

```
# SSH du returns no digit-leading line -> ParseDuSBOutput fails gracefully
podman machine ssh -- du -sb ... -> empty/garbage output -> ParseDuSBOutput error
```

## Preconditions

- Output with no parseable byte total must fail parsing.

## Steps

1. Set `req.Op` to `parse-du-invalid`.
2. Set `req.DuOutput` to garbage text.

## Context

- Distinguishes parse failure from graceful collection omission.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "parse-du-invalid"
	req.DuOutput = "du: cannot access '/missing': No such file or directory"
	return nil
}
```