const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

// Check initial state: no scanning-badge, no done-badge
const initScanning = await page.$$('[data-testid="scanning-badge"]');
const initDone = await page.$$('[data-testid="done-badge"]');
console.log(`INIT scanning-badge count: ${initScanning.length}`);
console.log(`INIT done-badge count: ${initDone.length}`);

// Click Start Scan
const startBtn = await page.$('[data-testid="start-scan-btn"]');
if (!startBtn) { console.log('BUTTON start-scan-btn: MISSING'); process.exitCode = 1; }
await startBtn.click();
await page.waitForTimeout(500);

// During scan: scanning-badge should appear on some cards
const midScanning = await page.$$('[data-testid="scanning-badge"]');
console.log(`MID scanning-badge count: ${midScanning.length}`);
for (const el of midScanning) {
    const visible = await el.isVisible();
    console.log(`MID scanning-badge visible: ${visible}`);
}

// Wait for scan to complete
await page.waitForTimeout(4000);

// After scan: done-badge should appear, scanning-badge should disappear
const finalScanning = await page.$$('[data-testid="scanning-badge"]');
const finalDone = await page.$$('[data-testid="done-badge"]');
console.log(`FINAL scanning-badge count: ${finalScanning.length}`);
console.log(`FINAL done-badge count: ${finalDone.length}`);

const anyDoneVisible = finalDone.length > 0;
console.log(`HAS_DONE_BADGES: ${anyDoneVisible}`);

console.log('DONE');
