const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

// Initial: no badges
const initScanning = await page.$$('[data-testid="scanning-badge"]');
const initDone = await page.$$('[data-testid="done-badge"]');
const initPending = await page.$$('[data-testid="pending-badge"]');
console.log(`INIT scanning-badge count: ${initScanning.length}`);
console.log(`INIT done-badge count: ${initDone.length}`);
console.log(`INIT pending-badge count: ${initPending.length}`);

// Click Start Scan
const startBtn = await page.$('[data-testid="start-scan-btn"]');
if (!startBtn) {
    console.log('BUTTON start-scan-btn: MISSING');
    console.log('DONE');
    process.exit(1);
}
await startBtn.click();

// Immediately after click: pending badges should appear on cards
await page.waitForTimeout(100);
const earlyPending = await page.$$('[data-testid="pending-badge"]');
console.log(`EARLY pending-badge count: ${earlyPending.length}`);

// Mid scan: some scanning badges appear, pending badges reduce
await page.waitForTimeout(2000);
const midScanning = await page.$$('[data-testid="scanning-badge"]');
const midPending = await page.$$('[data-testid="pending-badge"]');
const midDone = await page.$$('[data-testid="done-badge"]');
console.log(`MID scanning-badge count: ${midScanning.length}`);
console.log(`MID pending-badge count: ${midPending.length}`);
console.log(`MID done-badge count: ${midDone.length}`);

// Wait for scan to complete (start button reappears)
try {
    await page.waitForSelector('[data-testid="start-scan-btn"]', { state: 'visible', timeout: 45000 });
    console.log('SCAN_COMPLETE: start button reappeared');
} catch {
    console.log('SCAN_COMPLETE: timeout waiting for start button');
}

try {
    await page.waitForFunction(() => {
        return document.querySelectorAll('[data-testid="scanning-badge"]').length === 0
            && document.querySelectorAll('[data-testid="pending-badge"]').length === 0;
    }, { timeout: 3000 });
    console.log('BADGES_CLEARED: true');
} catch {
    console.log('BADGES_CLEARED: false');
}

await page.waitForTimeout(500);
const finalScanning = await page.$$('[data-testid="scanning-badge"]');
const finalPending = await page.$$('[data-testid="pending-badge"]');
const finalDone = await page.$$('[data-testid="done-badge"]');
console.log(`FINAL scanning-badge count: ${finalScanning.length}`);
console.log(`FINAL pending-badge count: ${finalPending.length}`);
console.log(`FINAL done-badge count: ${finalDone.length}`);

const hasDoneBadges = finalDone.length > 0;
console.log(`HAS_DONE_BADGES: ${hasDoneBadges}`);

console.log('DONE');
