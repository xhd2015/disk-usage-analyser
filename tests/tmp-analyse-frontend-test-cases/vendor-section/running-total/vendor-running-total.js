const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

const scanBtn = await page.$('[data-testid="vendor-scan-btn"]');
if (!scanBtn) {
    console.log('FAIL vendor-scan-btn: MISSING');
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

await scanBtn.click();

for (let i = 0; i < 120; i++) {
    await page.waitForTimeout(500);
    const done = await page.$('[data-testid="vendor-section"] [data-testid="vendor-done-badge"]');
    if (done) break;
}

const runningTotal = await page.$('[data-testid="vendor-running-total"]');
console.log(`ELEM vendor-running-total: ${runningTotal ? 'present' : 'MISSING'}`);

if (runningTotal) {
    const text = (await runningTotal.textContent()).trim();
    console.log(`TEXT vendor-running-total: "${text}"`);
    const hasSize = /\d/.test(text) && /[BKMGTP]/.test(text);
    console.log(`CHECK vendor-running-total-size: ${hasSize ? 'HAS_SIZE' : 'NO_SIZE'}`);
    if (!hasSize) process.exitCode = 1;
} else {
    process.exitCode = 1;
}

console.log('DONE');
