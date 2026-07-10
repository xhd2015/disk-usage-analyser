## Expected

- Root help is printed (no web server start).
- Subcommand list includes `analyse`, `scan`, `explain`, and `tmp-files`.
- `explain` is documented with a PATH argument (e.g. `explain [PATH]`).
- Server options `--dev` and `--component` appear in help text.
- Nested-help pointer: help mentions running a subcommand with `--help` (must contain
  both `command` and `--help`, e.g. `disk-usage-analyser <command> --help`).
- Root help does **not** list `inspect` as a standalone subcommand.
- User-facing help stdout ends with a trailing newline `\n`.

## Exit Code

- 0 (nil error from run)

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ServerWasStarted {
		t.Fatal("root -h must not start web server")
	}
	if resp.Err != nil {
		t.Fatalf("root -h must return nil error, got: %v", resp.Err)
	}
	// Prefer stdout; some flag libraries write help to stderr.
	out := resp.Stdout + resp.Stderr
	helpOut := resp.Stdout
	if helpOut == "" {
		helpOut = out
	}

	for _, want := range []string{
		"analyse",
		"scan",
		"explain",
		"tmp-files",
		"--dev",
		"--component",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("root help missing %q:\n%s", want, out)
		}
	}

	// Prefer the documented form "explain [PATH]" when present; also accept a line that
	// starts with explain and mentions PATH nearby.
	if !strings.Contains(out, "explain [PATH]") && !strings.Contains(out, "explain PATH") {
		found := false
		for _, line := range strings.Split(out, "\n") {
			trim := strings.TrimSpace(line)
			if strings.Contains(trim, "explain") && strings.Contains(strings.ToUpper(trim), "PATH") {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("root help should list explain with PATH (e.g. 'explain [PATH]'):\n%s", out)
		}
	}

	// Nested-help pointer: one line should tie "command" + "--help" together
	// (e.g. "disk-usage-analyser <command> --help for command-specific options").
	// Exact wording flexible; do not accept split tokens across unrelated lines only.
	hasNestedHelpPointer := false
	for _, line := range strings.Split(out, "\n") {
		ll := strings.ToLower(line)
		if strings.Contains(ll, "command") && strings.Contains(line, "--help") {
			hasNestedHelpPointer = true
			break
		}
	}
	if !hasNestedHelpPointer {
		t.Fatalf("root help must include a nested-help pointer line containing both \"command\" and \"--help\" (e.g. disk-usage-analyser <command> --help):\n%s", out)
	}

	// Reject standalone inspect subcommand lines like "inspect [FILE]" as a peer of scan.
	for _, line := range strings.Split(out, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "inspect ") || trim == "inspect" {
			t.Fatalf("root help must not list inspect as a subcommand; line %q\n%s", line, out)
		}
	}

	if !strings.HasSuffix(helpOut, "\n") {
		end := helpOut
		if len(end) > 24 {
			end = end[len(end)-24:]
		}
		if end == "" {
			end = "<empty>"
		}
		t.Fatalf("user-facing root help stdout must end with trailing \\n, got ending %q", end)
	}
}
```
