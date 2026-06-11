const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

// Initial: totals should be 0
const initTotal = await page.$eval('[data-testid="total-size"]', el => el.textContent);
const initReclaim = await page.$eval('[data-testid="reclaimable-size"]', el => el.textContent);
console.log(`INIT total: ${initTotal}`);
console.log(`INIT reclaimable: ${initReclaim}`);

// Click Start Scan
const startBtn = await page.$('[data-testid="start-scan-btn"]');
if (!startBtn) { console.log('BUTTON start-scan-btn: MISSING'); process.exitCode = 1; }
await startBtn.click();

// Mid scan: totals should be accumulating (not 0 Bytes)
await page.waitForTimeout(2000);
const midTotal = await page.$eval('[data-testid="total-size"]', el => el.textContent);
const midReclaim = await page.$eval('[data-testid="reclaimable-size"]', el => el.textContent);
console.log(`MID total: ${midTotal}`);
console.log(`MID reclaimable: ${midReclaim}`);

// Wait for scan to complete (start button reappears)
try {
    await page.waitForSelector('[data-testid="start-scan-btn"]', { state: 'visible', timeout: 15000 });
    console.log('SCAN_COMPLETE: start button reappeared');
} catch {
    console.log('SCAN_COMPLETE: timeout waiting for start button');
}

await page.waitForTimeout(500);
const finalTotal = await page.$eval('[data-testid="total-size"]', el => el.textContent);
const finalReclaim = await page.$eval('[data-testid="reclaimable-size"]', el => el.textContent);
console.log(`FINAL total: ${finalTotal}`);
console.log(`FINAL reclaimable: ${finalReclaim}`);

// Collect card sizes to verify total matches
const cardSizes = await page.$$eval('[data-testid="card-size"]', els => els.map(el => el.textContent));
console.log(`CARD_SIZES: ${JSON.stringify(cardSizes)}`);

console.log('DONE');
