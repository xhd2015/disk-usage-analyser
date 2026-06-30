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

const scanBtn = await page.$('[data-testid="worktrees-scan-btn"]');
await scanBtn.click();

for (let i = 0; i < 120; i++) {
    await page.waitForTimeout(500);
    const done = await page.$('[data-testid="worktrees-section"] [data-testid="worktrees-done-badge"]');
    if (done) break;
}

const repoSizes = await page.$$eval('[data-testid^="worktree-repo-size-"]', els =>
    els.map(el => el.textContent.trim())
).catch(() => []);

const repoBytes = repoSizes.map(parseHumanSize).filter(n => n > 0);
console.log(`COUNT worktree-repo-rows: ${repoSizes.length}`);
console.log(`REPO_SIZES_BYTES: ${JSON.stringify(repoBytes)}`);

const childGroups = await page.$$eval('[data-testid="worktrees-tree"] > div', groups =>
    groups.map(g => {
        const rows = g.querySelectorAll('[data-testid^="worktree-row-size-"]');
        return Array.from(rows).map(el => el.textContent.trim());
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

if (repoSizes.length < 2) {
    console.log('SKIP insufficient-worktree-repos');
    console.log('DONE');
    process.exit(0);
}

const repoSortOk = isMonotonicDesc(repoBytes);
console.log(`SORT worktree-repo-sizes: ${repoSortOk ? 'desc' : 'not-desc'}`);
console.log(`SORT worktree-child-sizes: ${childSortOk ? 'desc' : 'not-desc'}`);

if (repoSortOk && childSortOk) {
    console.log('CHECK worktree-sort-desc: pass');
} else {
    process.exitCode = 1;
}
console.log('DONE');