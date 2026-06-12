# Tmp Files Analyse — Frontend Test Cases

Run the tests:
```sh
doctest test -v ./
```

Feature: Tmp Files Analyse page — UI structure, navigation, scan interaction, stop behavior, and well-known software cache/log display.

## Test Tree

```
tmp-analyse-frontend-test-cases/
├── SETUP.md                              # Root: Request/Response + Run (calls playwright-debug)
├── verify-page-renders/                  # Page has heading, buttons, summary, core/software sections, collapse
│   ├── SETUP.md
│   ├── page-renders.js
│   └── ASSERT.md
├── verify-software-cards-render/         # Each software tool gets a card with label, size, reboot-safe badge
│   ├── SETUP.md
│   ├── software-cards-render.js
│   └── ASSERT.md
├── verify-not-detected-collapse/         # Non-detected tools grouped in collapsed "Not Detected" panel
│   ├── SETUP.md
│   ├── not-detected-collapse.js
│   └── ASSERT.md
├── verify-multi-path-breakdown/          # Go and Xcode cards show extra-path breakdown
│   ├── SETUP.md
│   ├── multi-path-breakdown.js
│   └── ASSERT.md
├── verify-cards-from-locations-event/    # Cards populated from SSE "locations" event before scan
│   ├── SETUP.md
│   ├── cards-from-locations.js
│   └── ASSERT.md
├── verify-navigation/                    # Nav link exists, click navigates to /tmp-analyse
│   ├── SETUP.md
│   ├── navigation.js
│   └── ASSERT.md
├── verify-scan-starts/                   # Click start, SSE events fire, cards update
│   ├── SETUP.md
│   ├── scan-starts.js
│   └── ASSERT.md
├── verify-stop-scan/                     # Click stop, scan halts, button reverts
│   ├── SETUP.md
│   ├── stop-scan.js
│   └── ASSERT.md
├── verify-pending-status/                # Cards show pending state before scan
├── verify-scan-progress/                 # Real-time size updates during scan
├── verify-scanning-indicator/            # Spinning icon shows while scanning
├── verify-totals-accumulate/             # Total/Reclaimable sizes accumulate during scan
└── verify-location-path-shown/           # Path is displayed after scan completes
```

## Test Cases

1. verify-page-renders — All expected DOM elements with data-testid exist, including core + software sections and collapse panel
2. verify-software-cards-render — Each of 17 software tools has a card with label, size, and reboot-safe badge
3. verify-not-detected-collapse — Non-detected tools appear in a collapsed "Not Detected" panel
4. verify-multi-path-breakdown — Go and Xcode cards show extra-path size breakdown (module vs build cache)
5. verify-cards-from-locations-event — All cards render from SSE "locations" event before scan button is clicked
6. verify-navigation — Nav link "Tmp Files" present, click goes to /tmp-analyse, page renders
7. verify-scan-starts — Click start triggers SSE events, card sizes update, button toggles
8. verify-stop-scan — Click stop reverts button, SSE stream closes
9. verify-pending-status — Cards show pending state before scan
10. verify-scan-progress — Real-time size updates during scan
11. verify-scanning-indicator — Spinning icon shows while scanning
12. verify-totals-accumulate — Total/Reclaimable sizes accumulate during scan
13. verify-location-path-shown — Path is displayed after scan completes

## Prerequisites

- Server must be running. Set `SERVER_URL` env var (default `http://localhost:8080`).
- Start with: `go run . --dev` or `go run .`

```sh
SERVER_URL=http://localhost:8080 doctest test tests/tmp-analyse-frontend-test-cases/
```
