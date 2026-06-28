# Scenario

**Feature**: Docker card shows runtime section after scan when stats available

```
# Docker location event includes runtimeItems -> runtime-section rendered
CollectRuntimeStats(docker) -> location event -> runtime-section / runtime-row-0
```

## Preconditions

- Docker card is detected (Docker Desktop data dir exists).
- When docker CLI returns stats, runtime-section appears with runtime-row-0.

## Steps

1. Set req.ScriptFile to docker-runtime-section.js.
2. Script completes scan and checks runtime-section or graceful absence.

## Context

- Test passes whether runtime stats are available or gracefully omitted.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "docker-runtime-section.js"
	return nil
}
```