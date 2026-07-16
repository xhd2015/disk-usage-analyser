# Scenario

**Decision**: full CLI subprocess path with real binary and real wall-clock idle shutdown

```
# harness builds module binary once per doctest session
go build -o <cache>/disk-usage-analyser .

# subprocess runs real CLI flags (no ServeForTest hook)
exec.Command(bin, --dev, --dev-idle-life <dur>) + NO_BROWSER=1
stdout -> parse "Serving directory preview at http://localhost:<port>"

# HTTP touch starts idle clock; real sleep waits for shutdown
GET /ping -> DevIdleWatch.Wrap -> Touch
time.Sleep past DevIdleLife -> shutdownDev() -> process exit, port closed
```

## Preconditions

- Depends on P1–P4 (`DevIdleWatch`, `--dev-idle-life` flag, `Serve` wiring, SSE touch).
- Nested root: `subprocess/DOCTEST.md` owns `Run`; parent `dev-idle/DOCTEST.md` does not run for these leaves.
- Module root resolved as `filepath.Join(DOCTEST_ROOT, "..", "..", "..")`.
- Session cache: `$TMPDIR/dev-idle-subprocess-doctest-<DOCTEST_SESSION_ID>/` with file lock;
  shared artifact is the compiled **`disk-usage-analyser`** binary (`go build -o`).
- Subprocess uses **`syscall.SysProcAttr{Setpgid: true}`** so teardown can signal the process group.
- **`NO_BROWSER=1`** in subprocess env (no browser open).
- Real wall-clock sleeps; account for default **10s** `DevIdleWatch` tick (≈ idle life + 10s);
  leaves labeled **`slow`** in ASSERT frontmatter.
- Port parsed from stdout line prefix **`Serving directory preview at http://localhost:`**.
- Teardown kills server PID if still running after assertions (cleanup only).
- `Run` holds **`subprocessRunMu`** so parallel leaves do not race on `FindAvailablePort` inside the binary.

## Steps

1. Build or reuse session-cached binary via `buildBinaryOnce`.
2. Start subprocess with `--dev` and leaf-specific `--dev-idle-life`.
3. Parse listening port from stdout; wait for TCP accept.
4. `GET /ping` once to start the idle clock.
5. `time.Sleep` past configured idle life (real wall clock).
6. Assert port closed, process exited, and stderr log per leaf.

```go
func Setup(t *testing.T, req *Request) error {
	req.BinPath = buildBinaryOnce(t)
	return nil
}
```