## Expected

- Named scan finds nested pnpm `node_modules` under `.pnpm/@babel+helpers@7.29.7/node_modules`.
- `named_enriched` for that path reports `packageManager=pnpm` (walk up to app-root `pnpm-lock.yaml`).
- Root `node_modules` `named_enriched` also reports `pnpm`.

## Errors

- Harness must not return error.

```go
import "strings"

func findHitBySuffix(hits []namedEventJSON, suffix string) (namedEventJSON, bool) {
	for _, hit := range hits {
		if strings.HasSuffix(hit.Path, suffix) {
			return hit, true
		}
	}
	return namedEventJSON{}, false
}

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	nestedSuffix := "node_modules/.pnpm/@babel+helpers@7.29.7/node_modules"
	nested, ok := findHitBySuffix(resp.NamedEnrichedJSON, nestedSuffix)
	if !ok {
		t.Fatalf("expected nested named_enriched hit ending with %q, got %d: %#v",
			nestedSuffix, len(resp.NamedEnrichedJSON), resp.NamedEnrichedJSON)
	}
	if nested.PackageManager != "pnpm" {
		t.Fatalf("nested named_enriched packageManager = %q, want pnpm (path=%q)", nested.PackageManager, nested.Path)
	}

	root, ok := findHitBySuffix(resp.NamedEnrichedJSON, "Projects/nested-pnpm-app/node_modules")
	if !ok {
		t.Fatalf("expected root named_enriched hit, got %#v", resp.NamedEnrichedJSON)
	}
	if root.PackageManager != "pnpm" {
		t.Fatalf("root named_enriched packageManager = %q, want pnpm", root.PackageManager)
	}
}
```