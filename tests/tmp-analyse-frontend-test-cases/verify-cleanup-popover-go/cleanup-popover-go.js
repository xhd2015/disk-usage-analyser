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

const indicator = await page.$('[data-testid="card-go"] [data-testid="cleanup-indicator"]');
if (!indicator) {
    console.log('ELEM go-cleanup-indicator: MISSING (go not detected)');
    console.log('SKIP go-popover-test');
    console.log('DONE');
    process.exit(0);
}

await indicator.click();
await page.waitForTimeout(500);

const popover = await page.$('[data-testid="cleanup-popover-go"]');
if (!popover) {
    console.log('ELEM cleanup-popover-go: MISSING');
    console.log('DONE');
    process.exit(1);
}

const popoverText = await popover.textContent() || '';

console.log(`CLEANUP_GO_CACHE: ${popoverText.includes('go clean -cache')}`);
console.log(`CLEANUP_GO_MODCACHE: ${popoverText.includes('go clean -modcache') || popoverText.includes('modcache')}`);
console.log(`CLEANUP_GO_RECOVERABLE: ${popoverText.includes('rebuilt') || popoverText.includes('Rebuilt') || popoverText.includes('re-downloaded') || popoverText.includes('Re-downloaded')}`);

console.log('DONE');
