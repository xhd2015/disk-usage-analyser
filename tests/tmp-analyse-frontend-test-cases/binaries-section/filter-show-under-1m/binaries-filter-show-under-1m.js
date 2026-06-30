const BASE = process.env.SERVER_URL || 'http://localhost:8080';
const ONE_MIB = 1048576;

function parseHumanSize(text) {
    const m = (text || '').trim().match(/^([\d.]+)\s*(Bytes|KB|MB|GB|TB)$/i);
    if (!m) return 0;
    const n = parseFloat(m[1]);
    const units = { bytes: 1, kb: 1024, mb: 1024 ** 2, gb: 1024 ** 3, tb: 1024 ** 4 };
    return Math.round(n * (units[m[2].toLowerCase()] || 1));
}

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

const filterBox = await page.$('[data-testid="binary-show-under-1m"]');
console.log(`ELEM binary-show-under-1m: ${filterBox ? 'present' : 'MISSING'}`);
if (!filterBox) {
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

const scanBtn = await page.$('[data-testid="binaries-scan-btn"]');
await scanBtn.click();

for (let i = 0; i < 120; i++) {
    await page.waitForTimeout(500);
    const done = await page.$('[data-testid="binaries-section"] [data-testid="binaries-done-badge"]');
    if (done) break;
}

const empty = await page.$('[data-testid="binaries-empty-state"]');
if (empty) {
    console.log('BINARIES_EMPTY_STATE: present');
    console.log('DONE');
    process.exit(0);
}

const countRows = async () => (await page.$$('[data-testid^="binary-row-"]')).length;

const countBefore = await countRows();
console.log(`COUNT binary-rows-before-filter: ${countBefore}`);

await filterBox.check();
await page.waitForTimeout(300);

const countAfter = await countRows();
console.log(`COUNT binary-rows-after-filter: ${countAfter}`);

// Detect whether any sub-1M binaries exist in full data (after showing all)
const sizeTexts = await page.$$eval('[data-testid^="binary-row-"] .ant-typography', els =>
    els.map(el => el.textContent.trim()).filter(t => /^\d/.test(t))
).catch(() => []);

let hasSub1M = false;
for (const text of sizeTexts) {
    const bytes = parseHumanSize(text);
    if (bytes > 0 && bytes < ONE_MIB) {
        hasSub1M = true;
        break;
    }
}

if (!hasSub1M) {
    console.log('SKIP no-sub-1m-binaries');
} else if (countAfter > countBefore) {
    console.log('CHECK binary-filter-toggle: pass');
} else {
    process.exitCode = 1;
}

console.log('DONE');