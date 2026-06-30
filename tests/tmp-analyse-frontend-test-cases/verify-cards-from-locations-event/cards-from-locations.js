const BASE = process.env.SERVER_URL || 'http://localhost:8080';

async function check(selector, label) {
    const el = await page.$(selector);
    if (!el) {
        console.log(`CARD_BEFORE_SCAN ${label}: MISSING`);
        return;
    }
    const visible = await el.isVisible();
    console.log(`CARD_BEFORE_SCAN ${label}: visible=${visible}`);
}

async function count(selector, label) {
    const els = await page.$$(selector);
    console.log(`COUNT_BEFORE_SCAN ${label}: ${els.length}`);
}

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

// Check core cards are visible before scan (from initial SSE locations event)
for (const cat of ['trash', 'temp', 'cache', 'log']) {
    await check(`[data-testid="card-${cat}"]`, cat);
}

// Count all cards present before scan
await count('[data-testid^="card-"]', 'all-cards');

// Count software cards specifically (non-core categories)
const softwareCats = ['go', 'npm', 'bun', 'yarn', 'pnpm', 'pip', 'cargo',
                      'ruby', 'docker', 'podman', 'nginx', 'gradle', 'maven',
                      'android', 'brew', 'xcode', 'composer'];
let swCount = 0;
for (const cat of softwareCats) {
    const card = await page.$(`[data-testid="card-${cat}"]`);
    if (card) {
        const visible = await card.isVisible();
        if (visible) swCount++;
    }
}
console.log(`SOFTWARE_CARDS_BEFORE_SCAN: count=${swCount}`);

console.log('DONE');
