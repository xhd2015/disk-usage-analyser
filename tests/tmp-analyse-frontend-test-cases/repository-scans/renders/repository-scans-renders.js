const BASE = process.env.SERVER_URL || 'http://localhost:8080';

async function check(selector, label) {
    const el = await page.$(selector);
    if (!el) {
        console.log(`ELEM ${label}: MISSING`);
        return false;
    }
    const visible = await el.isVisible();
    console.log(`ELEM ${label}: visible=${visible}`);
    return visible;
}

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

let ok = true;
ok = await check('[data-testid="section-repository-scans-heading"]', 'section-repository-scans-heading') && ok;
ok = await check('[data-testid="worktrees-section"]', 'worktrees-section') && ok;
ok = await check('[data-testid="worktrees-scan-btn"]', 'worktrees-scan-btn') && ok;
ok = await check('[data-testid="worktrees-stop-btn"]', 'worktrees-stop-btn') && ok;
ok = await check('[data-testid="binaries-section"]', 'binaries-section') && ok;
ok = await check('[data-testid="binaries-scan-btn"]', 'binaries-scan-btn') && ok;
ok = await check('[data-testid="binaries-stop-btn"]', 'binaries-stop-btn') && ok;

if (!ok) process.exitCode = 1;
console.log('DONE');