# Scenario

**Bug**: PM column text wraps ("unknow"/"n") when `unknown` exceeds narrow 56px track

## Steps

1. Set req.ScriptFile to table-columns-pm-no-wrap.js.
2. Run node_modules scan and inspect PM cell computed `white-space`.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "table-columns-pm-no-wrap.js"
	return nil
}
```