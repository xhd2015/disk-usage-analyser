## Expected Output

```
    size  symlinks  hardlinks  hardlink_size  shared_hardlink  shared_clone  pnpm_shared  bun_shared  path
  1024 B    0f+0d          0            0 B               0 B          0 B          0 B          0 B  .
```

## Expected

- Exit code 0.
- Stdout is two lines: aligned table header row, then summary for cwd.
- Summary `path` column uses `pathfmt.Short` (`.` when cwd is the analysed root).

## Errors

- No harness or CLI error.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/assert"
	"github.com/xhd2015/dot-pkgs/go-pkgs/pathfmt"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d (err=%v)", resp.ExitCode, resp.Err)
	}
	lines := strings.Split(strings.TrimSpace(resp.Stdout), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected header + summary, got %d lines:\n%s", len(lines), resp.Stdout)
	}
	for _, col := range []string{"size", "symlinks", "hardlinks", "hardlink_size", "shared_hardlink", "shared_clone", "pnpm_shared", "bun_shared", "path"} {
		if !strings.Contains(lines[0], col) {
			t.Fatalf("header missing column %q: %s", col, lines[0])
		}
	}
	if strings.Contains(lines[0], "\t") {
		t.Fatalf("expected aligned table header, got TSV: %s", lines[0])
	}
	shortPath := pathfmt.Short(req.FixtureDir)
	if shortPath != "." {
		t.Fatalf("pathfmt.Short(fixtureDir) with cwd in fixture: expected %q, got %q", ".", shortPath)
	}
	assert.Output(t, lines[1], ``+
`<contains>
0f+0d
1024 B
`+shortPath+`
</contains>`)
}
```