const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

// Check initial button state
const startBtn = await page.$('[data-testid="start-scan-btn"]');
const stopBtn = await page.$('[data-testid="stop-scan-btn"]');

if (!startBtn) { console.log('BUTTON start-scan-btn: MISSING'); process.exitCode = 1; }
const startVisible = startBtn ? await startBtn.isVisible() : false;
const stopVisible = stopBtn ? await stopBtn.isVisible() : false;
console.log(`BUTTON start-scan-btn: visible=${startVisible}`);
console.log(`BUTTON stop-scan-btn: visible=${stopVisible}`);

if (!startVisible) {
    console.log('FAIL: start-scan-btn should be visible initially');
    process.exitCode = 1;
}
if (stopVisible) {
    console.log('FAIL: stop-scan-btn should be hidden initially');
    process.exitCode = 1;
}

// Monitor SSE network requests
page.on('request', req => {
    if (req.url().includes('/api/tmp-analyse')) {
        console.log(`SSE_REQUEST: ${req.url()}`);
    }
});

// Click Start Scan
await startBtn.click();
await page.waitForTimeout(500);

// After click, start-scan-btn should hide, stop-scan-btn should show
try {
    const sv = await page.$eval('[data-testid="start-scan-btn"]', el => {
        return window.getComputedStyle(el).display !== 'none';
    });
    console.log(`BUTTON after click start-scan-btn: visible=${sv}`);
} catch {
    console.log('BUTTON after click start-scan-btn: visible=false');
}
try {
    const sv = await page.$eval('[data-testid="stop-scan-btn"]', el => {
        return window.getComputedStyle(el).display !== 'none';
    });
    console.log(`BUTTON after click stop-scan-btn: visible=${sv}`);
} catch {
    console.log('BUTTON after click stop-scan-btn: visible=false');
}

// Wait for SSE events and card updates
await page.waitForTimeout(3000);

// Check card sizes updated
const cardSizes = await page.$$eval('[data-testid="card-size"]', els =>
    els.map(el => el.textContent?.trim() || '')
);
for (const s of cardSizes) {
    console.log(`CARD_SIZE: "${s}"`);
}

// Check summary bar updated
const totalSize = await page.$eval('[data-testid="total-size"]', el =>
    el.textContent?.trim() || ''
).catch(() => '');
const reclaimableSize = await page.$eval('[data-testid="reclaimable-size"]', el =>
    el.textContent?.trim() || ''
).catch(() => '');
console.log(`SUMMARY total-size: "${totalSize}"`);
console.log(`SUMMARY reclaimable-size: "${reclaimableSize}"`);

console.log('DONE');
