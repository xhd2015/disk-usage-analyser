# Scenario

**Decision**: single `--name` value matches one or more directories

```
scan --name=node_modules -> dir match node_modules -> NamedHit
```

## Preconditions

- Exactly one `--name` value is specified.
- One or more directories with that basename exist under git repositories.

## Steps

1. Build fixtures with directories whose basename matches the single name.
2. Run scan with `--name <value>`.

## Context

- Each leaf in this subtree focuses on a distinct behavioral dimension (basic, nested, override-ignore, additive, json, size-accuracy, custom-root).

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{"scan"}
	}
	return nil
}
```
