# Scenario

**Decision**: `--name` flag — named entry detection additive to binary scanning

```
# --name flag adds basename-based directory/file hits alongside binary hits
scan --name=NAME -> scan_repo repo discovery -> repo walk -> name match -> NamedHit stream -> summary

# named dirs: recursive size computed, nested same-name dirs excluded from outer size
repo walk -> dir basename match --name -> computeDirSize(skipping nested same-name) -> NamedHit -> SkipDir

# named files: direct size, no SkipDir
repo walk -> file basename match --name -> NamedHit (direct size)

# non-matching entries: existing binary classification logic unchanged
repo walk -> non-matching entry -> detect/buildinfo -> BinaryHit
```

## Preconditions

- `--name` is a repeatable flag accepting a basename value.
- Named hits and binary hits are additive — both appear in the same scan output.
- Directory matches trigger recursive size computation with nested same-name exclusion.
- File matches report the file's direct size and do not skip directory traversal.
- Named-hunt operates within git repositories only, just like binary scanning.
- The `ignoredDirBasenames` map still applies for entries not matched by `--name`.

## Steps

1. Each leaf creates fixtures with specific name-match and/or binary scenarios.
2. Run `scan` with appropriate `--name` flags.
3. Assert NamedHit fields, stdout output, and summary line.

## Context

- The existing `ScanResult` gains a `NamedHits` field.
- Human output format: `<size>  name:<name>  <path>  (repo: <repoPath>)`
- JSON output adds `"type":"named"` to NDJSON objects.
- Summary format: `Found N binaries, M named entries, total <human-size>`

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{"scan"}
	}
	return nil
}

func writeData(t *testing.T, base string, rel string, data []byte) string {
	t.Helper()
	return writeFile(t, base, rel, data, 0644)
}

func mkNamedDir(t *testing.T, base string, rel string) string {
	t.Helper()
	return mkdir(t, base, rel)
}
```
