# Scenario

**Leaf**: human BREAKDOWN is an aligned multi-column table (SIZE / NAME / ROLE / RECLAIMABLE / NOTES)

```
# Unequal sizes force visible SIZE padding
explain PATH(android-avd: 400B,200B,100B,32B) ->
  BREAKDOWN
    SIZE   NAME  …  ROLE  RECLAIMABLE  NOTES
     400B  …           ☐  …
      32B  …           ☐  …
  # SIZE tokens right-aligned (same end column)
```

## Steps

1. Build AVD fixture (`MediumPhone.avd` with known unequal sizes).
2. Run `explain.RunCLI <avd-dir>` (default non-TTY → no color).

```go
func Setup(t *testing.T, req *Request) error {
	avdDir, _ := writeAVDFixture(t, req.FixtureDir)
	req.TargetPath = avdDir
	req.Args = []string{avdDir}
	return nil
}
```
