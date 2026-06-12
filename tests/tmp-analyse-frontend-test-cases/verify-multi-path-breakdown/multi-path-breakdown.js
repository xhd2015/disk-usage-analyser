const BASE = process.env.SERVER_URL || 'http://localhost:8080';

async function check(selector, label) {
    const el = await page.$(selector);
    if (!el) {
        console.log(`ELEM ${label}: MISSING`);
        return null;
    }
    const text = (await el.textContent()) || '';
    const visible = await el.isVisible();
    console.log(`ELEM ${label}: "${text.trim()}" visible=${visible}`);
    return el;
}

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

// Check Go card for unified breakdown
await check('[data-testid="card-go"] [data-testid="breakdown-items"]', 'card-go-breakdown-items');
await check('[data-testid="card-go"] [data-testid="breakdown-label-1"]', 'card-go-breakdown-label-1');
await check('[data-testid="card-go"] [data-testid="breakdown-size-1"]', 'card-go-breakdown-size-1');
await check('[data-testid="card-go"] [data-testid="breakdown-row-1"]', 'card-go-breakdown-row-1');

// Check Xcode card for unified breakdown
await check('[data-testid="card-xcode"] [data-testid="breakdown-items"]', 'card-xcode-breakdown-items');
await check('[data-testid="card-xcode"] [data-testid="breakdown-label-1"]', 'card-xcode-breakdown-label-1');
await check('[data-testid="card-xcode"] [data-testid="breakdown-size-1"]', 'card-xcode-breakdown-size-1');
await check('[data-testid="card-xcode"] [data-testid="breakdown-row-1"]', 'card-xcode-breakdown-row-1');

// Verify breakdown labels show full tilde paths (not truncated)
const goLabelText = await page.$eval('[data-testid="card-go"] [data-testid="breakdown-label-1"]', el => el.textContent).catch(() => '');
console.log(`FULL_PATH go-label-starts-with-tilde: ${goLabelText.startsWith('~/')}`);
console.log(`FULL_PATH go-label-not-truncated: ${goLabelText.split('/').length > 3}`);

const xcodeLabelText = await page.$eval('[data-testid="card-xcode"] [data-testid="breakdown-label-1"]', el => el.textContent).catch(() => '');
console.log(`FULL_PATH xcode-label-starts-with-tilde: ${xcodeLabelText.startsWith('~/')}`);
console.log(`FULL_PATH xcode-label-not-truncated: ${xcodeLabelText.split('/').length > 3}`);

// Single-path tools should NOT have breakdown-items
const singlePathCats = ['npm', 'bun', 'docker', 'gradle', 'maven'];
for (const cat of singlePathCats) {
    const breakdownEl = await page.$(`[data-testid="card-${cat}"] [data-testid="breakdown-items"]`);
    if (breakdownEl) {
        console.log(`ELEM card-${cat}-breakdown-items: HAS_UNEXPECTED_BREAKDOWN`);
    } else {
        console.log(`ELEM card-${cat}-breakdown-items: MISSING (expected, single-path)`);
    }
}

console.log('DONE');
