# Scenario

**Leaf**: `-v` prints a warning when CloudStorage paths are skipped during scan

## Preconditions

- Same fixture as `skip-cloud-storage`.

## Steps

1. Build the CloudStorage and local fixture repos.
2. Run `scan -v`.

## Context

- Warning policy matches `scan_repo`: only emit remote-backed skip warnings with `-v`.

```go
func Setup(t *testing.T, req *Request) error {
	provider := cloudStorageProvider(t, req.HomeDir, "GoogleDrive-user@example.com")
	mkdir(t, provider, ".shortcut-targets-by-id")
	cloudRepo := repo(t, req.HomeDir, "Library/CloudStorage/GoogleDrive-user@example.com/Projects/cloud-app")
	writeMachO(t, cloudRepo, "bin/cloud-app")
	localRepo := repo(t, req.HomeDir, "Projects/local-app")
	writeMachO(t, localRepo, "bin/local-app")
	req.Args = []string{"scan", "-v"}
	return nil
}
```