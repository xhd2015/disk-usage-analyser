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

// Section headings
await check('[data-testid="section-core-heading"]', 'section-core-heading');
await check('[data-testid="section-software-heading"]', 'section-software-heading');

// Core category cards
await check('[data-testid="card-trash"]', 'card-trash');
await check('[data-testid="card-temp"]', 'card-temp');
await check('[data-testid="card-cache"]', 'card-cache');
await check('[data-testid="card-log"]', 'card-log');
await check('[data-testid="card-swap"]', 'card-swap');

// Card internals for core categories
for (const cat of ['trash', 'temp', 'cache', 'log', 'swap']) {
    await check(`[data-testid="card-${cat}"] [data-testid="card-label"]`, `card-${cat}-label`);
    await check(`[data-testid="card-${cat}"] [data-testid="card-size"]`, `card-${cat}-size`);
    await check(`[data-testid="card-${cat}"] [data-testid="reboot-safe-badge"]`, `card-${cat}-reboot-safe`);
}

// Software cards section exists
await check('[data-testid="section-software"]', 'section-software');

// Count all cards (core + software)
await count('[data-testid^="card-"]', 'all-cards');

// Not Detected collapse section
await check('[data-testid="collapse-not-detected"]', 'collapse-not-detected');

// Check some well-known software cards
await check('[data-testid="card-go"]', 'card-go');
await check('[data-testid="card-npm"]', 'card-npm');
await check('[data-testid="card-docker"]', 'card-docker');
await check('[data-testid="card-bun"]', 'card-bun');

console.log('DONE');
