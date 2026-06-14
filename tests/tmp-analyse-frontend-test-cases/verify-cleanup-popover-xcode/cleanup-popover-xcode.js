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

const indicator = await page.$('[data-testid="card-xcode"] [data-testid="cleanup-indicator"]');
if (!indicator) {
    console.log('ELEM xcode-cleanup-indicator: MISSING (xcode not detected)');
    console.log('SKIP xcode-popover-test');
    console.log('DONE');
    process.exit(0);
}

await indicator.click();
await page.waitForTimeout(500);

const popover = await page.$('[data-testid="cleanup-popover-xcode"]');
if (!popover) {
    console.log('ELEM cleanup-popover-xcode: MISSING');
    console.log('DONE');
    process.exit(1);
}

const popoverText = await popover.textContent() || '';

console.log(`CLEANUP_XCODE_DERIVED: ${popoverText.includes('DerivedData')}`);
console.log(`CLEANUP_XCODE_SIMULATORS: ${popoverText.includes('simctl')}`);
console.log(`CLEANUP_XCODE_RECOVERABLE: ${popoverText.includes('rebuilt') || popoverText.includes('Rebuilt') || popoverText.includes('recreated') || popoverText.includes('Recreated')}`);

console.log('DONE');
