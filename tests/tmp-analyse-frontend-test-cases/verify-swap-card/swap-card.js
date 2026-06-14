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

await check('[data-testid="card-swap"]', 'card-swap');
await check('[data-testid="card-swap"] [data-testid="card-label"]', 'card-swap-label');
await check('[data-testid="card-swap"] [data-testid="card-size"]', 'card-swap-size');
await check('[data-testid="card-swap"] [data-testid="reboot-safe-badge"]', 'card-swap-reboot-safe');

const labelEl = await page.$('[data-testid="card-swap"] [data-testid="card-label"]');
if (labelEl) {
    const labelText = (await labelEl.textContent()) || '';
    console.log(`SWAP_LABEL_CHECK: ${labelText.includes('Swap')}`);
}

console.log('DONE');
