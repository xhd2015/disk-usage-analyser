# Tmp Analyse Test Cases

Run the tests:
```sh
doctest test -v ./
```

Feature: Tmp Files Analyse — analyze temp files, caches, and logs from the system
and popular developer tools. Shows how much space can be freed.

## Test Tree

```
tmp-analyse-test-cases/
├── SETUP.md                              # Root: Request/Response types
├── verify-discover-locations/            # DiscoverLocations returns 22+ macOS paths with Detected flags
├── verify-software-locations/            # Each software location has correct Path, Label, Category, ExtraPaths
├── verify-detected-by-existence/         # Detected flag set via os.Stat; core always true, software conditional
├── verify-multi-path-locations/          # Go and Xcode have ExtraPaths; single-path tools have none
├── verify-initial-locations-event/       # First SSE event is "locations" with full 22+ item array
├── verify-unsupported-platform/          # Non-darwin handler sends unsupported_platform event
├── verify-extra-path-scan/               # Multi-path locations accumulate ExtraSizes/ExtraCounts
├── verify-locations-rest-endpoint/       # REST endpoint returns locations JSON without scanning
├── verify-calculate-size/                # CalculateSize sums file sizes correctly
├── verify-empty-dir/                     # CalculateSize handles empty directories
├── verify-nested-dirs/                   # CalculateSize recursively sums nested dirs
├── verify-summary-totals/                # BuildSummary computes total + reclaimable
├── verify-sse-format/                    # SSE handler emits all expected event types
├── verify-progress-stream/               # ScanWithProgress fires progress callbacks
├── verify-scan-with-partial-error/       # Error handling during scan
└── verify-totals-in-progress/            # BuildProgressPayload accumulates totals correctly
```

## Test Cases

1. verify-discover-locations — DiscoverLocations returns 22+ locations (5 core + 17 software) with Detected flags
2. verify-software-locations — Each software location has correct label, category, RebootSafe=true; Go and Xcode have ExtraPaths
3. verify-detected-by-existence — Real existing dir (e.g. /tmp) gets Detected=true; non-existing gets Detected=false
4. verify-multi-path-locations — Go has 2 ExtraPaths, Xcode has 1; single-path tools have none
5. verify-initial-locations-event — Handler sends "event: locations" as first SSE event with full JSON array
6. verify-unsupported-platform — Non-darwin sends unsupported_platform event and no scan
7. verify-extra-path-scan — Multi-path scanning populates ExtraSizes/ExtraCounts arrays
8. verify-locations-rest-endpoint — REST endpoint returns locations JSON array without scanning
9. verify-calculate-size — CalculateSize on mock FS with 3 files returns correct total bytes and file count
9. verify-empty-dir — CalculateSize on empty FS returns size=0, count=0, no error
10. verify-nested-dirs — CalculateSize recursively sums files in nested subdirectories
11. verify-summary-totals — BuildSummary correctly separates total vs reclaimable (rebootSafe) sizes
12. verify-sse-format — SSE handler emits locations, location, summary, done events with valid JSON
13. verify-progress-stream — ScanWithProgress fires progress callbacks
14. verify-scan-with-partial-error — Error handling during scan
15. verify-totals-in-progress — BuildProgressPayload accumulates totals correctly
```
