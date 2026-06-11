const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

// Record initial card sizes
const initialSizes = await page.$$eval('[data-testid="card-size"]', els =>
    els.map(el => el.textContent?.trim() || '')
);
for (const s of initialSizes) {
    console.log(`INITIAL_SIZE: "${s}"`);
}

// Monitor SSE
page.on('request', req => {
    if (req.url().includes('/api/tmp-analyse')) {
        console.log(`SSE_REQUEST: ${req.url()}`);
    }
});

// Click Start Scan
const startBtn = await page.$('[data-testid="start-scan-btn"]');
if (!startBtn) { console.log('BUTTON start-scan-btn: MISSING'); process.exitCode = 1; }
await startBtn.click();

// Wait for scanning to progress
await page.waitForTimeout(500);

// Capture intermediate sizes
const midSizes = await page.$$eval('[data-testid="card-size"]', els =>
    els.map(el => el.textContent?.trim() || '')
);
for (const s of midSizes) {
    console.log(`MID_SIZE: "${s}"`);
}

// Wait for scan to complete
await page.waitForTimeout(4000);

// Capture final sizes
const finalSizes = await page.$$eval('[data-testid="card-size"]', els =>
    els.map(el => el.textContent?.trim() || '')
);
for (const s of finalSizes) {
    console.log(`FINAL_SIZE: "${s}"`);
}

// Verify some sizes changed from initial 0
const hasNonZero = finalSizes.some(s => s !== '0 Bytes');
console.log(`HAS_NONZERO: ${hasNonZero}`);

console.log('DONE');
