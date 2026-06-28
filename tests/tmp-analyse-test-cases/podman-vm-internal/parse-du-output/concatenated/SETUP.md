# Scenario

**Feature**: concatenated `du -sb` line (no tab) parses to bytes

```
# SSH du returns "BYTESPATH" without tab -> ParseDuSBOutput still extracts BYTES
podman machine ssh -- du -sb /home/core/.../storage -> ParseDuSBOutput -> 1735229168
```

## Preconditions

- Line `1735229168/home/core/.local/share/containers/storage` parses to 1735229168 bytes.

## Steps

1. Set `req.Op` to `parse-du-concat`.
2. Set `req.DuOutput` from fixture.

## Context

- Observed on some Podman machine SSH sessions where tab is omitted.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "parse-du-concat"
	req.FixtureFile = "du-concatenated.txt"
	return nil
}
```