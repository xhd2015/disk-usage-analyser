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

// === Go Card (multi-item: 2 items) ===

// Check Go card exists
const goCard = await page.$('[data-testid="card-go"]');
if (goCard) {
    console.log('CARD_EXISTS go: true');

    // Must have unified breakdown-items wrapper
    await check('[data-testid="card-go"] [data-testid="breakdown-items"]', 'go-breakdown-items');

    // Row 0: primary path (go/pkg/mod) with its size
    await check('[data-testid="card-go"] [data-testid="breakdown-row-0"]', 'go-breakdown-row-0');
    await check('[data-testid="card-go"] [data-testid="breakdown-label-0"]', 'go-breakdown-label-0');
    await check('[data-testid="card-go"] [data-testid="breakdown-size-0"]', 'go-breakdown-size-0');

    // Row 1: extra path (Caches/go-build) with its size
    await check('[data-testid="card-go"] [data-testid="breakdown-row-1"]', 'go-breakdown-row-1');
    await check('[data-testid="card-go"] [data-testid="breakdown-label-1"]', 'go-breakdown-label-1');
    await check('[data-testid="card-go"] [data-testid="breakdown-size-1"]', 'go-breakdown-size-1');

    // Must NOT have standalone card-path (all paths are in breakdown)
    const goCardPath = await page.$('[data-testid="card-go"] [data-testid="card-path"]');
    console.log(`NO_STANDALONE_PATH go: ${goCardPath === null}`);

    // Row layout: both rows should use flexbox
    const goRow0 = await page.$('[data-testid="card-go"] [data-testid="breakdown-row-0"]');
    if (goRow0) {
        const display = await goRow0.evaluate(el => getComputedStyle(el).display);
        const justify = await goRow0.evaluate(el => getComputedStyle(el).justifyContent);
        console.log(`ROW_LAYOUT go-row0-display: ${display}`);
        console.log(`ROW_LAYOUT go-row0-justify: ${justify}`);
    }
} else {
    console.log('CARD_EXISTS go: false (skipping breakdown checks)');
}

// === Xcode Card (multi-item: 5 items, macOS only) ===
const xcodeCard = await page.$('[data-testid="card-xcode"]');
if (xcodeCard) {
    console.log('CARD_EXISTS xcode: true');

    await check('[data-testid="card-xcode"] [data-testid="breakdown-items"]', 'xcode-breakdown-items');

    for (let i = 0; i <= 4; i++) {
        await check(`[data-testid="card-xcode"] [data-testid="breakdown-row-${i}"]`, `xcode-breakdown-row-${i}`);
        await check(`[data-testid="card-xcode"] [data-testid="breakdown-label-${i}"]`, `xcode-breakdown-label-${i}`);
        await check(`[data-testid="card-xcode"] [data-testid="breakdown-size-${i}"]`, `xcode-breakdown-size-${i}`);
    }

    const xcodeCardPath = await page.$('[data-testid="card-xcode"] [data-testid="card-path"]');
    console.log(`NO_STANDALONE_PATH xcode: ${xcodeCardPath === null}`);

    const xcodeRow0 = await page.$('[data-testid="card-xcode"] [data-testid="breakdown-row-0"]');
    if (xcodeRow0) {
        const display = await xcodeRow0.evaluate(el => getComputedStyle(el).display);
        const justify = await xcodeRow0.evaluate(el => getComputedStyle(el).justifyContent);
        console.log(`ROW_LAYOUT xcode-row0-display: ${display}`);
        console.log(`ROW_LAYOUT xcode-row0-justify: ${justify}`);
    }
} else {
    console.log('CARD_EXISTS xcode: false (skipping breakdown checks)');
}

// === Single-item cards: must still use card-path centered (NO breakdown-items) ===
for (const cat of ['npm', 'bun', 'docker', 'gradle']) {
    const card = await page.$(`[data-testid="card-${cat}"]`);
    if (!card) {
        console.log(`CARD_EXISTS ${cat}: false (skipping)`);
        continue;
    }

    // Must have card-path
    const cpEl = await page.$(`[data-testid="card-${cat}"] [data-testid="card-path"]`);
    console.log(`HAS_CARD_PATH ${cat}: ${cpEl !== null}`);

    // Must NOT have breakdown-items
    const bdEl = await page.$(`[data-testid="card-${cat}"] [data-testid="breakdown-items"]`);
    console.log(`NO_BREAKDOWN_ITEMS ${cat}: ${bdEl === null}`);
}

console.log('DONE');
