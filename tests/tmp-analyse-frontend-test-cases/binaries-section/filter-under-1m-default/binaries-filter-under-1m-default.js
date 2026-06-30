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
if (filterBox) {
    const checked = await filterBox.isChecked();
    console.log(`CHECKBOX binary-show-under-1m: ${checked ? 'checked' : 'unchecked'}`);
}

const scanBtn = await page.$('[data-testid="binaries-scan-btn"]');
if (!scanBtn) {
    console.log('FAIL binaries-scan-btn: MISSING');
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

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

const sizeTexts = await page.$$eval('[data-testid^="binary-row-"] .ant-typography', els =>
    els.map(el => el.textContent.trim()).filter(t => /^\d/.test(t))
).catch(() => []);

let under1mVisible = false;
for (const text of sizeTexts) {
    const bytes = parseHumanSize(text);
    console.log(`BINARY_ROW_SIZE: "${text}" -> ${bytes}`);
    if (bytes > 0 && bytes < ONE_MIB) {
        under1mVisible = true;
    }
}

console.log(`UNDER_1M_VISIBLE: ${under1mVisible ? 'yes' : 'no'}`);
if (!filterBox || under1mVisible) {
    process.exitCode = 1;
} else {
    console.log('CHECK binary-filter-default: pass');
}
console.log('DONE');