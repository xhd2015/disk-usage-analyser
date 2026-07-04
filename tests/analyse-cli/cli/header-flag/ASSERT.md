## Expected Output

```
    size  symlinks  hardlinks  hardlink_size  shared_hardlink  shared_clone  pnpm_shared  bun_shared  path
     0 B    0f+0d          0            0 B               0 B          0 B          0 B          0 B  <pathfmt.Short(fixture)>
```

## Expected

- Exit code 0.
- First stdout line is the aligned table header with all column names (always printed).
- Second line is the zero-byte summary row with `pathfmt.Short` in the path column.
- `--header` is accepted even though the header row is always emitted.

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
	assert.Output(t, lines[1], ``+
`<contains>
0 B
0f+0d
`+shortPath+`
</contains>`)
}
```