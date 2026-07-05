# Scenario

**Leaf**: vendor scan skips pass 2 — no `named_enriched`, quick `scan_complete` then `done`

## Steps

1. Create git repo with `vendor/`.
2. Run `named-scan` with `name=vendor`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Name = "vendor"
	vendorRepo(t, req.HomeDir, "Projects/vendor-app")
	return nil
}
```