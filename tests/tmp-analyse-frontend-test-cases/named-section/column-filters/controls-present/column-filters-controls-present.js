const BASE = process.env.SERVER_URL || 'http://localhost:8080';

async function readFilterValue(testId) {
    const el = await page.$(`[data-testid="${testId}"]`);
    if (!el) return null;

    const tag = await el.evaluate(node => node.tagName.toLowerCase());
    if (tag === 'select') {
        return el.evaluate(node => node.value);
    }

    const selected = await el.$('.ant-segmented-item-selected, .ant-select-selection-item');
    if (selected) {
        return selected.evaluate(node => (node.textContent || '').trim().toLowerCase());
    }

    return el.evaluate(node => (node.getAttribute('data-value') || node.textContent || '').trim().toLowerCase());
}

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

const scanBtn = await page.$('[data-testid="node-modules-scan-btn"]');
if (!scanBtn) {
    console.log('FAIL node-modules-scan-btn: MISSING');
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

await scanBtn.click();

for (let i = 0; i < 120; i++) {
    await page.waitForTimeout(500);
    const done = await page.$('[data-testid="node-modules-section"] [data-testid="node-modules-done-badge"]');
    if (done) break;
}

const empty = await page.$('[data-testid="node-modules-empty-state"]');
if (empty) {
    console.log('NAMED_EMPTY_STATE: present');
    console.log('DONE');
    process.exit(0);
}

const gitFilter = await page.$('[data-testid="node-modules-filter-git"]');
const pkgFilter = await page.$('[data-testid="node-modules-filter-package-json"]');
const pmFilter = await page.$('[data-testid="node-modules-filter-pm"]');

console.log(`ELEM node-modules-filter-git: ${gitFilter ? 'present' : 'MISSING'}`);
console.log(`ELEM node-modules-filter-package-json: ${pkgFilter ? 'present' : 'MISSING'}`);
console.log(`ELEM node-modules-filter-pm: ${pmFilter ? 'present' : 'MISSING'}`);

if (!gitFilter || !pkgFilter || !pmFilter) {
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

const gitVal = await readFilterValue('node-modules-filter-git');
const pkgVal = await readFilterValue('node-modules-filter-package-json');
const pmVal = await readFilterValue('node-modules-filter-pm');

console.log(`FILTER git: ${gitVal}`);
console.log(`FILTER packageJson: ${pkgVal}`);
console.log(`FILTER pm: ${pmVal}`);

const norm = v => {
    const s = (v || '').toLowerCase();
    if (s === 'all' || s.includes('all')) return 'all';
    return s;
};

const gitNorm = norm(gitVal);
const pkgNorm = norm(pkgVal);
const pmNorm = norm(pmVal);

if (gitNorm === 'all' && pkgNorm === 'all' && pmNorm === 'all') {
    console.log('FILTER_DEFAULT: git=all packageJson=all pm=all');
} else {
    process.exitCode = 1;
}

console.log('DONE');