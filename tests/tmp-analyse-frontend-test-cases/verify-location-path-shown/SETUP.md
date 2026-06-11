## Preconditions
- Each location card shows its actual file path after scan completes
- Path element has data-testid="card-path" and contains a non-empty value

## Steps
1. Set req.ScriptFile to "location-path.js"
2. The script starts a scan, waits for completion, checks that paths are shown

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "location-path.js"
	return nil
}
```
