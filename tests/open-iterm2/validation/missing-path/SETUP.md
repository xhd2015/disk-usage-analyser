# Scenario

**Leaf**: empty JSON body without `path` returns HTTP 400

## Steps

1. POST `{}` to `/api/open-iterm2`.

```go
func Setup(t *testing.T, req *Request) error {
	req.BodyJSON = "{}"
	return nil
}
```