const BASE = process.env.SERVER_URL || 'http://localhost:8080';

async function check(selector, label) {
    const el = await page.$(selector);
    if (!el) {
        console.log(`ELEM ${label}: MISSING`);
        return null;
    }
    const text = (await el.textContent()) || '';
    console.log(`ELEM ${label}: "${text.trim()}"`);
    return el;
}

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

const npmCard = await page.$('[data-testid="card-npm"]');
if (!npmCard) {
    console.log('SKIP npm-breakdown: npm card not found (npm not detected)');
    console.log('DONE');
    process.exit(0);
}

const startBtn = await page.$('[data-testid="start-scan-btn"]');
if (!startBtn) {
    console.log('SKIP npm-breakdown: start scan button not found');
    console.log('DONE');
    process.exit(0);
}

await startBtn.click();
await page.waitForTimeout(3000);

const breakdownItems = await page.$('[data-testid="card-npm"] [data-testid="breakdown-items"]');
if (breakdownItems) {
    console.log(`NPM_HAS_BREAKDOWN: true`);
    const rows = await page.$$('[data-testid="card-npm"] [data-testid^="breakdown-row-"]');
    console.log(`NPM_BREAKDOWN_COUNT: ${rows.length}`);
    for (let i = 0; i < rows.length; i++) {
        const labelEl = await page.$(`[data-testid="card-npm"] [data-testid="breakdown-label-${i}"]`);
        const sizeEl = await page.$(`[data-testid="card-npm"] [data-testid="breakdown-size-${i}"]`);
        if (labelEl) {
            const labelText = (await labelEl.textContent()) || '';
            console.log(`NPM_BREAKDOWN_LABEL_${i}: "${labelText.trim()}"`);
        }
        if (sizeEl) {
            const sizeText = (await sizeEl.textContent()) || '';
            console.log(`NPM_BREAKDOWN_SIZE_${i}: "${sizeText.trim()}"`);
        }
    }
} else {
    console.log(`NPM_HAS_BREAKDOWN: false`);
    const singlePath = await page.$('[data-testid="card-npm"] [data-testid="card-path"]');
    if (singlePath) {
        console.log(`NPM_SINGLE_PATH: true`);
    } else {
        console.log(`NPM_SINGLE_PATH: false`);
    }
}

console.log('DONE');
