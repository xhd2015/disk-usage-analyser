const BASE = process.env.SERVER_URL || 'http://localhost:8080';

function parseHumanSize(text) {
    const m = (text || '').trim().match(/^([\d.]+)\s*(Bytes|KB|MB|GB|TB)$/i);
    if (!m) return 0;
    const n = parseFloat(m[1]);
    const units = { bytes: 1, kb: 1024, mb: 1024 ** 2, gb: 1024 ** 3, tb: 1024 ** 4 };
    return Math.round(n * (units[m[2].toLowerCase()] || 1));
}

function isMonotonicDesc(nums) {
    for (let i = 1; i < nums.length; i++) {
        if (nums[i] > nums[i - 1]) return false;
    }
    return true;
}

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

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

const repoTotals = await page.$$eval('[data-testid="binaries-tree"] > div', groups =>
    groups.map(g => {
        const sizeEl = g.querySelector('.ant-typography');
        return sizeEl ? sizeEl.textContent.trim() : '';
    }).filter(Boolean)
).catch(() => []);

const repoBytes = repoTotals.map(parseHumanSize).filter(n => n > 0);
console.log(`REPO_TOTALS_BYTES: ${JSON.stringify(repoBytes)}`);

const childGroups = await page.$$eval('[data-testid="binaries-tree"] > div', groups =>
    groups.map(g => {
        const rows = g.querySelectorAll('[data-testid^="binary-row-"]');
        return Array.from(rows).map(row => {
            const texts = row.querySelectorAll('.ant-typography');
            const sizeText = texts.length ? texts[texts.length - 1].textContent.trim() : '';
            return sizeText;
        }).filter(Boolean);
    })
).catch(() => []);

let childSortOk = true;
for (const sizes of childGroups) {
    const bytes = sizes.map(parseHumanSize).filter(n => n > 0);
    if (bytes.length > 1 && !isMonotonicDesc(bytes)) {
        childSortOk = false;
        console.log(`CHILD_SIZES_BYTES: ${JSON.stringify(bytes)}`);
    }
}

const totalRows = childGroups.reduce((n, g) => n + g.length, 0);
console.log(`COUNT binary-rows: ${totalRows}`);

if (totalRows < 2) {
    console.log('SKIP insufficient-binary-rows');
    console.log('DONE');
    process.exit(0);
}

const repoSortOk = repoBytes.length < 2 || isMonotonicDesc(repoBytes);
console.log(`SORT binary-repo-totals: ${repoSortOk ? 'desc' : 'not-desc'}`);
console.log(`SORT binary-child-sizes: ${childSortOk ? 'desc' : 'not-desc'}`);

if (repoSortOk && childSortOk) {
    console.log('CHECK binary-sort-desc: pass');
} else {
    process.exitCode = 1;
}
console.log('DONE');