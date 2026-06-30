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

await scanBtn.click();

for (let i = 0; i < 120; i++) {
    await page.waitForTimeout(500);
    const done = await page.$('[data-testid="worktrees-section"] [data-testid="worktrees-done-badge"]');
    if (done) break;
}

const tree = await page.$('[data-testid="worktrees-tree"]');
console.log(`ELEM worktrees-tree: ${tree ? 'present' : 'MISSING'}`);

const repoRows = await page.$$('[data-testid^="worktree-repo-row-"]');
console.log(`COUNT worktree-repo-rows: ${repoRows.length}`);

const wtRows = await page.$$('[data-testid^="worktree-row-"]');
console.log(`COUNT worktree-rows: ${wtRows.length}`);

if (repoRows.length > 0) {
    const sizeEl = await page.$('[data-testid^="worktree-repo-size-"]');
    const sizeText = sizeEl ? (await sizeEl.textContent()).trim() : '';
    console.log(`WORKTREE_REPO_SIZE: "${sizeText}"`);
}

if (wtRows.length > 0) {
    const firstSize = await page.$eval('[data-testid^="worktree-row-size-"]', el => el.textContent.trim()).catch(() => '');
    console.log(`WORKTREE_ROW_SIZE: "${firstSize}"`);
}

const mainBadge = await page.$('[data-testid="worktree-main-badge"]');
console.log(`ELEM worktree-main-badge: ${mainBadge ? 'present' : 'ABSENT'}`);

if (!tree || repoRows.length === 0 || mainBadge) process.exitCode = 1;
console.log('DONE');