# Scenario

**Leaf**: live `--top` ranking matches inspect of an equivalent capture tree

## Steps

1. Write on-disk tree matching `sampleInspectTree`: `huge.bin` (400), `mid.bin` (200), `big/deep.bin` (50).
2. Write equivalent TreeResult JSON (`min` field) with the same absolute paths/sizes beside the fixture.
3. Run live `RunCLI --min 1B --top 2 <fixture>`.
4. Assert re-runs `RunCLI --inspect <json> --top 2` and compares match order (`huge.bin`, `mid.bin`).

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "huge.bin", 400)
	writeSizedFile(t, req.FixtureDir, "mid.bin", 200)
	writeSizedFile(t, req.FixtureDir, "big/deep.bin", 50)

	jsonPath := filepath.Join(filepath.Dir(req.FixtureDir), "capture.json")
	tree, total := sampleInspectTree(req.FixtureDir)
	writeTreeResultJSON(t, jsonPath, req.FixtureDir, total, 0, 24, tree)

	req.Args = []string{"--min", "1B", "--top", "2", req.FixtureDir}
	return nil
}
```
