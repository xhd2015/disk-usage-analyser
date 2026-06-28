# Scenario

**Feature**: page-renders UI verification

```
# React page renders tmp-analyse cards and scan UI
User -> /tmp-analyse -> playwright-debug checks data-testid elements
```

## Preconditions
- The React app is compiled and served at `SERVER_URL`
- No `/tmp-analyse` page exists yet, or a skeleton page exists

## Steps
1. Set req.ScriptFile to "page-renders.js"
2. The Run function executes playwright-debug with the script
3. The script navigates to /tmp-analyse and checks all data-testid elements

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "page-renders.js"
	return nil
}
```
