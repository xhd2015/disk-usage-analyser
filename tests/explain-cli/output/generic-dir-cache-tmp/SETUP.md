# Scenario

**Leaf**: generic-dir basename remap — `Cache` / `tmp` → reclaimable roles with `☑`, size DESC

```
# Fixture (generic-dir, not specialized kind):
#   Cache/entry  200B → role cache  ☑
#   tmp/work     100B → role tmp    ☑
#   notes.txt     32B → role file   ☐
explain PATH(generic-dir) ->
  KIND: generic-dir
  BREAKDOWN size DESC: Cache, tmp, notes.txt
  roles cache/tmp reclaimable; notes neutral
```

## Steps

1. Build generic cache/tmp fixture under `req.FixtureDir`.
2. Run `explain.RunCLI <fixture-dir>`.

```go
func Setup(t *testing.T, req *Request) error {
	target := writeGenericCacheTmpFixture(t, req.FixtureDir)
	req.TargetPath = target
	req.Args = []string{target}
	return nil
}
```
