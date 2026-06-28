# Scenario

**Feature**: tab-separated `du -sb` line parses to bytes

```
# SSH du returns "BYTES\tPATH" -> ParseDuSBOutput extracts BYTES
podman machine ssh -- du -sb /home/core/.../storage -> ParseDuSBOutput -> 1735229168
```

## Preconditions

- Line `1735229168\t/home/core/.local/share/containers/storage` parses to 1735229168 bytes.

## Steps

1. Set `req.Op` to `parse-du-tab`.
2. Set `req.DuOutput` from fixture.

## Context

- Standard `du -sb` format with tab separator.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "parse-du-tab"
	req.FixtureFile = "du-tab-separated.txt"
	return nil
}
```