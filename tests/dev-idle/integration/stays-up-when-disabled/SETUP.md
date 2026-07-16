# Scenario

**Leaf**: dev server keeps listening when `DevIdleLife` is zero (disabled)

```
DevIdleLife=0 -> DevIdleWatch disabled -> no idle shutdown
GET /ping -> advance clock 5s -> port still listening
```

## Steps

1. Set `DevIdleLife` to 0 (disabled).
2. Start `ServeForTest` with `Dev: true`.
3. `GET /ping` once.
4. Advance fake clock 5s and pump idle checks.
5. Assert port still accepts connections; stop server in harness cleanup.

```go
func Setup(t *testing.T, req *Request) error {
	req.Scenario = "stays-up-when-disabled"
	req.DevIdleLife = 0
	return nil
}
```