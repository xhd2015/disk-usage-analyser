# Scenario

**Decision**: request validation errors

```
POST /api/open-iterm2 -> missing or invalid path -> HTTP 400
```

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "open-iterm2"
	return nil
}
```