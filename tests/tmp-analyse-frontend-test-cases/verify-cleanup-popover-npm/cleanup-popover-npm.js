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

const indicator = await page.$('[data-testid="card-npm"] [data-testid="cleanup-indicator"]');
if (!indicator) {
    console.log('ELEM npm-cleanup-indicator: MISSING (npm not detected)');
    console.log('SKIP npm-popover-test');
    console.log('DONE');
    process.exit(0);
}

await indicator.click();
await page.waitForTimeout(500);

const popover = await page.$('[data-testid="cleanup-popover-npm"]');
if (!popover) {
    console.log('ELEM cleanup-popover-npm: MISSING');
    console.log('DONE');
    process.exit(1);
}

const popoverText = await popover.textContent() || '';

console.log(`CLEANUP_NPM_CACHE_CLEAN: ${popoverText.includes('npm cache clean')}`);
console.log(`CLEANUP_NPM_CACHE_VERIFY: ${popoverText.includes('npm cache verify')}`);
console.log(`CLEANUP_NPM_RECOVERABLE: ${popoverText.includes('re-downloaded') || popoverText.includes('Re-downloaded')}`);
console.log(`CLEANUP_NPM_HAS_DESCRIPTION: ${popoverText.length > 30}`);

console.log('DONE');
