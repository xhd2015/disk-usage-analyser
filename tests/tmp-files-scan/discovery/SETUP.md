# Scenario

**Decision**: repository discovery scope

```
discovery -> repo membership | ignore dirs | max depth | permission denied | remote fs | no repos
```

## Preconditions

- Only files under discovered git repositories may be reported.
- Default ignored basenames include `.git`, `vendor`, `node_modules`, and `.venv`.

## Steps

1. Build fixture directory trees that isolate one discovery rule per leaf.
2. Run `scan` with deterministic fake roots.

## Context

- Permission errors must not abort the whole scan when other repos can be scanned.
- Remote-backed cloud storage paths must warn on stderr and be skipped before walking.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{"scan"}
	}
	return nil
}
```
