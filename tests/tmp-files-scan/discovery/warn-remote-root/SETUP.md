# Scenario

**Leaf**: an explicit `--root` under CloudStorage is warned and skipped before walking

## Preconditions

- Users may pass `--root ~/Library/CloudStorage/<provider>` directly.
- Remote-backed roots should emit a warning and avoid walking into them.

## Steps

1. Create `~/Library/CloudStorage/GoogleDrive-user@example.com/.shortcut-targets-by-id`.
2. Create a git repo with a Mach-O binary inside that CloudStorage tree.
3. Run `scan --root` pointing at the CloudStorage provider directory.

## Context

- Covers the case where the scan root itself is a remote-backed filesystem.

```go
func Setup(t *testing.T, req *Request) error {
	provider := cloudStorageProvider(t, req.HomeDir, "GoogleDrive-user@example.com")
	mkdir(t, provider, ".shortcut-targets-by-id")
	cloudRepo := repo(t, req.HomeDir, "Library/CloudStorage/GoogleDrive-user@example.com/Projects/cloud-app")
	writeMachO(t, cloudRepo, "bin/cloud-app")
	req.Args = []string{
		"scan",
		"--root",
		"~/Library/CloudStorage/GoogleDrive-user@example.com",
	}
	return nil
}
```