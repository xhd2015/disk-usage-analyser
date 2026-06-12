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

// Check collapse panel exists
await check('[data-testid="collapse-not-detected"]', 'collapse-not-detected');

// Check header
await check('[data-testid="collapse-not-detected"] .ant-collapse-header', 'not-detected-header');

// Click to expand and check items inside
const collapseHeader = await page.$('[data-testid="collapse-not-detected"] .ant-collapse-header');
if (collapseHeader) {
    await collapseHeader.click();
    await page.waitForTimeout(300);

    const notDetectedItems = await page.$$('[data-testid="collapse-not-detected"] [data-testid="not-detected-item"]');
    for (const item of notDetectedItems) {
        const text = (await item.textContent()) || '';
        console.log(`NOT_DETECTED_ITEM: "${text.trim()}"`);
    }

    const notDetectedList = await page.$$('[data-testid="collapse-not-detected"] [data-testid="not-detected-item-name"]');
    for (const item of notDetectedList) {
        const text = (await item.textContent()) || '';
        console.log(`NOT_DETECTED_NAME: "${text.trim()}"`);
    }
}

console.log('DONE');
