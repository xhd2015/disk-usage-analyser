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

async function count(selector, label) {
    const els = await page.$$(selector);
    console.log(`COUNT ${label}: ${els.length}`);
}

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

// Upper left: page heading
await check('[data-testid="page-heading"]', 'page-heading');

// Upper right: buttons
await check('[data-testid="start-scan-btn"]', 'start-scan-btn');
await check('[data-testid="stop-scan-btn"]', 'stop-scan-btn');

// Summary bar
await check('[data-testid="summary-bar"]', 'summary-bar');
await check('[data-testid="total-size"]', 'total-size');
await check('[data-testid="reclaimable-size"]', 'reclaimable-size');

// Category cards
await check('[data-testid="card-trash"]', 'card-trash');
await check('[data-testid="card-temp"]', 'card-temp');
await check('[data-testid="card-cache"]', 'card-cache');
await check('[data-testid="card-log"]', 'card-log');

// Card internals
await count('[data-testid^="card-"]', 'card-');
for (const cat of ['trash', 'temp', 'cache', 'log']) {
    await check(`[data-testid="card-${cat}"] [data-testid="card-label"]`, `card-${cat}-label`);
    await check(`[data-testid="card-${cat}"] [data-testid="card-size"]`, `card-${cat}-size`);
    await check(`[data-testid="card-${cat}"] [data-testid="reboot-safe-badge"]`, `card-${cat}-reboot-safe`);
}

console.log('DONE');
