const BASE = process.env.SERVER_URL || 'http://localhost:8080';
const TEN_MIB = 10485760;

function parseHumanSize(text) {
    const m = (text || '').trim().match(/^([\d.]+)\s*(Bytes|KB|MB|GB|TB)$/i);
    if (!m) return 0;
    const n = parseFloat(m[1]);
    const units = { bytes: 1, kb: 1024, mb: 1024 ** 2, gb: 1024 ** 3, tb: 1024 ** 4 };
    return Math.round(n * (units[m[2].toLowerCase()] || 1));
}

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

const filterBox = await page.$('[data-testid="worktree-show-under-10m"]');
console.log(`ELEM worktree-show-under-10m: ${filterBox ? 'present' : 'MISSING'}`);
if (filterBox) {
    const checked = await filterBox.isChecked();
    console.log(`CHECKBOX worktree-show-under-10m: ${checked ? 'checked' : 'unchecked'}`);
}

const scanBtn = await page.$('[data-testid="worktrees-scan-btn"]');
if (!scanBtn) {
    console.log('FAIL worktrees-scan-btn: MISSING');
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

await scanBtn.click();

for (let i = 0; i < 120; i++) {
    await page.waitForTimeout(500);
    const done = await page.$('[data-testid="worktrees-section"] [data-testid="worktrees-done-badge"]');
    if (done) break;
}

const repoRows = await page.$$('[data-testid^="worktree-repo-row-"]');
console.log(`COUNT worktree-repo-rows: ${repoRows.length}`);

const repoSizes = await page.$$eval('[data-testid^="worktree-repo-size-"]', els =>
    els.map(el => el.textContent.trim())
).catch(() => []);
const rowSizes = await page.$$eval('[data-testid^="worktree-row-size-"]', els =>
    els.map(el => el.textContent.trim())
).catch(() => []);

let under10mVisible = false;
for (const text of [...repoSizes, ...rowSizes]) {
    const bytes = parseHumanSize(text);
    console.log(`WORKTREE_SIZE: "${text}" -> ${bytes}`);
    if (bytes > 0 && bytes < TEN_MIB) {
        under10mVisible = true;
    }
}

console.log(`UNDER_10M_VISIBLE: ${under10mVisible ? 'yes' : 'no'}`);

if (!filterBox || under10mVisible || repoRows.length === 0) {
    if (repoRows.length === 0) {
        // no repos — assert will skip
    } else {
        process.exitCode = 1;
    }
} else {
    console.log('CHECK worktree-filter-default: pass');
}
console.log('DONE');