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

const xcodeCard = await page.$('[data-testid="card-xcode"]');
if (!xcodeCard) {
    console.log('SKIP xcode-breakdown-five: Xcode card not detected');
    console.log('DONE');
    process.exit(0);
}

await check('[data-testid="card-xcode"] [data-testid="breakdown-items"]', 'xcode-breakdown-items');

for (let i = 0; i <= 4; i++) {
    await check(`[data-testid="card-xcode"] [data-testid="breakdown-row-${i}"]`, `xcode-breakdown-row-${i}`);
    await check(`[data-testid="card-xcode"] [data-testid="breakdown-label-${i}"]`, `xcode-breakdown-label-${i}`);
    await check(`[data-testid="card-xcode"] [data-testid="breakdown-size-${i}"]`, `xcode-breakdown-size-${i}`);
}

const label4 = await page.$eval('[data-testid="card-xcode"] [data-testid="breakdown-label-4"]', el => el.textContent).catch(() => '');
console.log(`FULL_PATH xcode-row4-documentation-cache: ${label4.includes('DocumentationCache')}`);

console.log('DONE');