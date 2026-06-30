# Scenario

**Leaf**: cloud-synced Library/CloudStorage paths are detected before walking and skipped with a warning

## Preconditions

- macOS Google Drive and similar providers live under `~/Library/CloudStorage/<provider>/`.
- `.shortcut-targets-by-id` is a known remote-backed marker that can hang on `fdopendir`.
- A local repository outside CloudStorage must still be scanned.

## Steps

1. Create `~/Library/CloudStorage/GoogleDrive-user@example.com/.shortcut-targets-by-id`.
2. Create a git repo with a Mach-O binary inside that CloudStorage tree.
3. Create a readable local git repo at `~/Projects/local-app` with a Mach-O binary.
4. Run `scan` from the fixture home.

## Context

- Reproduces the hang reported when scanning `~` reaches Google Drive paths.
- Expected fix: warn on stderr and skip the remote-backed directory before walking it.

```go
func Setup(t *testing.T, req *Request) error {
	provider := cloudStorageProvider(t, req.HomeDir, "GoogleDrive-user@example.com")
	mkdir(t, provider, ".shortcut-targets-by-id")
	cloudRepo := repo(t, req.HomeDir, "Library/CloudStorage/GoogleDrive-user@example.com/Projects/cloud-app")
	writeMachO(t, cloudRepo, "bin/cloud-app")
	localRepo := repo(t, req.HomeDir, "Projects/local-app")
	writeMachO(t, localRepo, "bin/local-app")
	req.Args = []string{"scan"}
	return nil
}
```