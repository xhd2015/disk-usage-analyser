# Scenario

**Feature**: package manager label from project root lockfile markers

```
named hit node_modules -> inspect parent dir -> first lockfile match -> packageManager on SSE named event
```

## Preconditions

- Detection inspects the parent directory of each `node_modules` path.
- Priority: `bun.lockb`/`bun.lock` → `pnpm-lock.yaml` → `package-lock.json` → `yarn.lock` → `node_modules/.pnpm/` → parse `package.json` `packageManager` field → `package.json` exists (default `npm`) → `unknown`.
- Each leaf uses an isolated fake-home repo with a single `node_modules` hit.

## Context

- Sibling leaves are mutually exclusive on the winning signal (lockfile, `packageManager` field, or `package.json` presence).
- `has-package-json` / `no-package-json` leaves assert the `hasPackageJson` SSE field.
- `package-json-default-npm` and `package-json-package-manager-*` leaves cover Corepack and default-npm paths without lockfiles.
- `lockfile-wins-over-field` asserts lockfile priority over the `packageManager` field.
- `nested-pnpm-node-modules` asserts nested `.pnpm/pkg@ver/node_modules` paths resolve PM from the app-root lockfile.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "named-scan"
	req.Name = "node_modules"
	return nil
}
```