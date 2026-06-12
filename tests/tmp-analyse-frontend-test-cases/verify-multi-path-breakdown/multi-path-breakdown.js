const BASE = process.env.SERVER_URL || 'http://localhost:8080';

async function check(selector, label) {
    const el = await page.$(selector);
    if (!el) {
        console.log(`ELEM ${label}: MISSING`);
        return;
    }
    const text = (await el.textContent()) || '';
    const visible = await el.isVisible();
    console.log(`ELEM ${label}: "${text.trim()}" visible=${visible}`);
}

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

// Check Go card for extra breakdown
await check('[data-testid="card-go"] [data-testid="extra-breakdown"]', 'card-go-extra-breakdown');
await check('[data-testid="card-go"] [data-testid="extra-breakdown-label-0"]', 'card-go-extra-label-0');
await check('[data-testid="card-go"] [data-testid="extra-breakdown-size-0"]', 'card-go-extra-size-0');

// Check Xcode card for extra breakdown
await check('[data-testid="card-xcode"] [data-testid="extra-breakdown"]', 'card-xcode-extra-breakdown');
await check('[data-testid="card-xcode"] [data-testid="extra-breakdown-label-0"]', 'card-xcode-extra-label-0');
await check('[data-testid="card-xcode"] [data-testid="extra-breakdown-size-0"]', 'card-xcode-extra-size-0');

// Single-path tools should NOT have extra breakdown
const singlePathCats = ['npm', 'bun', 'docker', 'gradle', 'maven'];
for (const cat of singlePathCats) {
    const breakdownEl = await page.$(`[data-testid="card-${cat}"] [data-testid="extra-breakdown"]`);
    if (breakdownEl) {
        console.log(`ELEM card-${cat}-extra-breakdown: HAS_UNEXPECTED_BREAKDOWN`);
    } else {
        console.log(`ELEM card-${cat}-extra-breakdown: MISSING (expected, single-path)`);
    }
}

console.log('DONE');
