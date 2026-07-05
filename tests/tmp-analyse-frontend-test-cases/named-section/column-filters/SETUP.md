# Scenario

**Feature**: node_modules table column filters for Git, package.json, and package manager

```
# size filter then column filters -> visible repo-grouped rows
filterNamedRepos(showUnder1M) -> filterNamedReposByColumnFilters(git, pkgjson, pm) -> node_modules tree rows

# UI controls in column header row drive client filter state
node-modules-filter-git / package-json / pm -> column filter state -> row visibility
```

## Preconditions

- Column filters apply after the existing `<1M` size filter on each `NamedHit`.
- Git and package.json filters are tri-state (`all` / `yes` / `no`); PM filter is `all` or a specific manager.
- Repos with zero visible hits after column filters are omitted.
- Pure logic leaves live under `filter-named-repos-by-column/` (nested doctest root).
- UI leaves use Playwright against the dev server started by the parent `Run`.

## Steps

1. Pure logic: run `column-filters-harness.ts` via nested `filter-named-repos-by-column` doctest root.
2. UI automation: set `req.ScriptFile` per leaf and drive filter controls after scan.

```go
func Setup(t *testing.T, req *Request) error {
	_ = req.ScriptFile
	return nil
}
```