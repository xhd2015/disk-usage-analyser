# Tmp Analyse Test Cases

Feature: Tmp Files Analyse — analyze temp files and directories on the system,
including trash size, to see how much space can be manually freed before reboot.

## Test Tree

```
tmp-analyse-test-cases/
├── SETUP.md                              # Root: Request/Response types
├── verify-discover-locations/            # DiscoverLocations returns correct macOS paths
├── verify-calculate-size/                # CalculateSize sums file sizes correctly
├── verify-empty-dir/                     # CalculateSize handles empty directories
├── verify-nested-dirs/                   # CalculateSize recursively sums nested dirs
├── verify-summary-totals/                # BuildSummary computes total + reclaimable
└── verify-sse-format/                    # SSE handler emits valid event stream
```

## Test Cases

1. verify-discover-locations — DiscoverLocations("/Users/testuser") returns paths including ~/.Trash, ~/Library/Caches, ~/Library/Logs with correct categories
2. verify-calculate-size — CalculateSize on a mock FS with 3 files returns correct total bytes and file count
3. verify-empty-dir — CalculateSize on an empty FS returns size=0, count=0, no error
4. verify-nested-dirs — CalculateSize recursively sums files in nested subdirectories
5. verify-summary-totals — BuildSummary correctly separates total vs reclaimable (rebootSafe) sizes
6. verify-sse-format — SSE handler emits properly formatted events with valid JSON payloads
