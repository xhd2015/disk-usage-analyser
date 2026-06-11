# Tmp Files Analyse — Frontend Test Cases

Feature: Tmp Files Analyse page — UI structure, navigation, scan interaction, and stop behavior.

## Test Tree

```
tmp-analyse-frontend-test-cases/
├── SETUP.md                              # Root: Request/Response + Run (calls playwright-debug)
├── verify-page-renders/                  # Page has heading, buttons, summary bar, 4 category cards
│   ├── SETUP.md
│   ├── page-renders.js
│   └── ASSERT.md
├── verify-navigation/                    # Nav link exists, click navigates to /tmp-analyse
│   ├── SETUP.md
│   ├── navigation.js
│   └── ASSERT.md
├── verify-scan-starts/                   # Click start, SSE events fire, cards update
│   ├── SETUP.md
│   ├── scan-starts.js
│   └── ASSERT.md
└── verify-stop-scan/                     # Click stop, scan halts, button reverts
    ├── SETUP.md
    ├── stop-scan.js
    └── ASSERT.md
```

## Test Cases

1. verify-page-renders — All expected DOM elements with data-testid exist with correct initial state
2. verify-navigation — Nav link "Tmp Files" present, click goes to /tmp-analyse, page renders
3. verify-scan-starts — Click start triggers SSE events, card sizes update, button toggles
4. verify-stop-scan — Click stop reverts button, SSE stream closes

## Prerequisites

- Server must be running. Set `SERVER_URL` env var (default `http://localhost:8080`).
- Start with: `go run . --dev` or `go run .`

```sh
SERVER_URL=http://localhost:8080 doctest test tests/tmp-analyse-frontend-test-cases/
```
