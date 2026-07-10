# Scenario

**Leaf**: root `-h` lists subcommands, documents server flags, and points to nested help (no server start)

```
# Root help (no subcommand)
disk-usage-analyser -h|--help
  -> lists analyse | scan | explain | tmp-files
  -> documents --dev and --component
  -> points users to: disk-usage-analyser <command> --help
  -> exit 0, StartServer never called
```

## Steps

1. Call `run.RunWithOptions` with `-h` (root help, not `explain -h`).

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "dispatch"
	req.Args = []string{"-h"}
	return nil
}
```
