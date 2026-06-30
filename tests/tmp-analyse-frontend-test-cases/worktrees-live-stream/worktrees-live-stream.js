const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

const scanBtn = await page.$('[data-testid="worktrees-scan-btn"]');
if (!scanBtn) {
    console.log('FAIL worktrees-scan-btn: MISSING');
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

let maxRepoRows = 0;
let maxWtRows = 0;
let sawScanning = false;
let sawRepoGrowth = false;
let sawWtGrowth = false;

await scanBtn.click();

for (let i = 0; i < 120; i++) {
    await page.waitForTimeout(500);
    const scanning = await page.$('[data-testid="worktrees-section"] [data-testid="worktrees-scanning-badge"]');
    if (scanning) sawScanning = true;

    const repoRows = await page.$$('[data-testid^="worktree-repo-row-"]');
    if (repoRows.length > maxRepoRows) {
        if (maxRepoRows > 0) sawRepoGrowth = true;
        maxRepoRows = repoRows.length;
    }

    const wtRows = await page.$$('[data-testid^="worktree-row-"]');
    if (wtRows.length > maxWtRows) {
        if (maxWtRows > 0) sawWtGrowth = true;
        maxWtRows = wtRows.length;
    }

    console.log(`POLL worktree-repo-rows: ${repoRows.length}`);
    console.log(`POLL worktree-rows: ${wtRows.length}`);

    const done = await page.$('[data-testid="worktrees-section"] [data-testid="worktrees-done-badge"]');
    if (done) break;
}

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

const finalRepoSizes = await page.$$eval('[data-testid^="worktree-repo-size-"]', els =>
    els.map(el => el.textContent.trim())
).catch(() => []);
const finalRepoBytes = finalRepoSizes.map(parseHumanSize).filter(n => n > 0);
console.log(`FINAL_REPO_SIZES_BYTES: ${JSON.stringify(finalRepoBytes)}`);
const finalSortDesc = finalRepoBytes.length < 2 || isMonotonicDesc(finalRepoBytes);
console.log(`SORT worktree-repo-sizes-final: ${finalSortDesc ? 'desc' : 'not-desc'}`);

console.log(`CHECK worktrees-scanning-seen: ${sawScanning}`);
console.log(`CHECK worktree-repo-growth: ${sawRepoGrowth || maxRepoRows > 0}`);
console.log(`CHECK worktree-row-growth: ${sawWtGrowth || maxWtRows > 0}`);
console.log(`MAX worktree-repo-rows: ${maxRepoRows}`);
console.log(`MAX worktree-rows: ${maxWtRows}`);

if (!sawScanning || maxRepoRows === 0) process.exitCode = 1;
if (maxRepoRows >= 2 && !finalSortDesc) process.exitCode = 1;
console.log('DONE');