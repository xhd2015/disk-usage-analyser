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
if (!filterBox) {
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

const scanBtn = await page.$('[data-testid="worktrees-scan-btn"]');
await scanBtn.click();

for (let i = 0; i < 120; i++) {
    await page.waitForTimeout(500);
    const done = await page.$('[data-testid="worktrees-section"] [data-testid="worktrees-done-badge"]');
    if (done) break;
}

const countVisible = async () => {
    const repos = (await page.$$('[data-testid^="worktree-repo-row-"]')).length;
    const rows = (await page.$$('[data-testid^="worktree-row-"]')).length;
    return repos + rows;
};

const before = await countVisible();
console.log(`COUNT worktree-visible-before: ${before}`);
console.log(`COUNT worktree-repo-rows-before: ${(await page.$$('[data-testid^="worktree-repo-row-"]')).length}`);

await filterBox.check();
await page.waitForTimeout(300);

const after = await countVisible();
console.log(`COUNT worktree-visible-after: ${after}`);

const allSizes = await page.$$eval(
    '[data-testid^="worktree-repo-size-"], [data-testid^="worktree-row-size-"]',
    els => els.map(el => el.textContent.trim())
).catch(() => []);

let hasSub10M = false;
for (const text of allSizes) {
    const bytes = parseHumanSize(text);
    if (bytes > 0 && bytes < TEN_MIB) {
        hasSub10M = true;
        break;
    }
}

if (!hasSub10M) {
    console.log('SKIP no-sub-10m-worktrees');
} else if (after > before) {
    console.log('CHECK worktree-filter-toggle: pass');
} else {
    process.exitCode = 1;
}

console.log('DONE');