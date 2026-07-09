# Scenario

**Feature**: live FS scan with query flags (shared View on LiveTreeSource)

```
scan PATH --min 1B --top N -> LiveTreeSource -> View(tree + matches)
```

## Preconditions

- Live defaults still apply unless overridden (`--min 1M` would hide small fixtures).
- Leaves use `--min 1B` (or similar) so small temp files participate in rankings.
- Match ranking uses the full live tree; tree section still follows live max-depth defaults.

## Context

- Option B applies equally to live and inspect sources.
- Parity leaf compares live top ranking to inspect of an equivalent capture JSON.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	return nil
}
```
