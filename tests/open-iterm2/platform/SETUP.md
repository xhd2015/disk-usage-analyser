# Scenario

**Decision**: platform guard for non-macOS

```
POST /api/open-iterm2 -> KOOL_ITERM2_GOOS=linux -> HTTP 501
```

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "open-iterm2"
	req.GoOS = "linux"
	return nil
}
```