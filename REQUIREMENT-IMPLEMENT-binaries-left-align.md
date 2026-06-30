# REQUIREMENT-IMPLEMENT — Binary files left align

## Context

Followup: binaries-tree not left-aligned (computed text-align: center). Worktrees section already fixed.

## Tests sealed — do not modify

```
tests/tmp-analyse-frontend-test-cases/binaries-section/left-aligned/
tests/tmp-analyse-frontend-test-cases/binaries-section/SETUP.md (if staged)
tests/tmp-analyse-frontend-test-cases/DOCTEST.md (index only)
```

## Implementation

`disk-usage-analyser-react/src/TmpFilesAnalyse.tsx`:

Apply same styling as `worktrees-tree` to `binaries-tree`:
- `style={{ width: '100%', textAlign: 'left' }}` on `binaries-tree` div
- `textAlign: 'left'` on repo row containers and child indent div

Rebuild: `cd disk-usage-analyser-react && bun run build`

## Verify

```sh
doctest test ./tests/tmp-analyse-frontend-test-cases/binaries-section/left-aligned
doctest test ./tests/tmp-analyse-frontend-test-cases/binaries-section/...
doctest test ./...
```