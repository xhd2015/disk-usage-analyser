# Scenario

**Leaf**: explain on main_*.sqlite inside SeaTalk Application Support prefers seatalk-app-support parent context

## Steps

1. Build SeaTalk fixture (`Application Support/SeaTalk`).
2. Run `explain.RunCLI` with the absolute path to `main_1.sqlite`.

```go
func Setup(t *testing.T, req *Request) error {
	_, mainDB := writeSeaTalkFixture(t, req.FixtureDir)
	req.TargetPath = mainDB
	req.Args = []string{mainDB}
	return nil
}
```
