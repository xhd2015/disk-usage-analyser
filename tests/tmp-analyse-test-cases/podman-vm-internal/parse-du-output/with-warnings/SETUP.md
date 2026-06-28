# Scenario

**Feature**: noisy `du -sb` output with stderr warnings still parses

```
# SSH du mixes permission warnings with valid line -> ParseDuSBOutput picks digit-leading line
podman machine ssh -- du -sb ... -> noisy stdout -> ParseDuSBOutput -> 1735229168
```

## Preconditions

- stderr permission warnings mixed with a valid `du` line still yield 1735229168 bytes.

## Steps

1. Set `req.Op` to `parse-du-warnings`.
2. Load noisy fixture containing warnings plus valid line.

## Context

- Parser prefers last matching line or first line starting with digits.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "parse-du-warnings"
	req.FixtureFile = "du-with-warnings.txt"
	return nil
}
```