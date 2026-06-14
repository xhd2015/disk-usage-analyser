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

const swapCard = await page.$('[data-testid="card-swap"]');
if (!swapCard) {
    console.log('SKIP swap-non-reclaimable: swap card not found');
    console.log('DONE');
    process.exit(0);
}

const nonReclaimableEl = await page.$('[data-testid="card-swap"] [data-testid="non-reclaimable-badge"]');
if (nonReclaimableEl) {
    const text = (await nonReclaimableEl.textContent()) || '';
    console.log(`SWAP_NON_RECLAIMABLE_EXISTS: true`);
    console.log(`SWAP_NON_RECLAIMABLE_TEXT: "${text.trim()}"`);
} else {
    console.log('ELEM non-reclaimable-badge: MISSING');
}

console.log('DONE');
