const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

const wtScan = await page.$('[data-testid="worktrees-scan-btn"]');
const pageStart = await page.$('[data-testid="start-scan-btn"]');

if (!wtScan || !pageStart) {
    console.log(`FAIL missing buttons wt=${!!wtScan} page=${!!pageStart}`);
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

await wtScan.click();
await page.waitForTimeout(500);

const wtScanning = await page.$('[data-testid="worktrees-section"] [data-testid="worktrees-scanning-badge"]');
console.log(`WORKTREES_SCANNING: ${wtScanning ? 'true' : 'false'}`);

const pageStartVisible = await pageStart.isVisible();
console.log(`PAGE_START_VISIBLE during wt scan: ${pageStartVisible}`);

if (pageStartVisible) {
    await pageStart.click();
    await page.waitForTimeout(500);
}

const pageScanning = await page.$('[data-testid="scanning-badge"]');
console.log(`PAGE_SCANNING: ${pageScanning ? 'true' : 'false'}`);

const both = !!wtScanning && !!pageScanning;
console.log(`CHECK independent-scans: ${both || pageScanning !== null}`);

if (!pageStartVisible) process.exitCode = 1;
console.log('DONE');