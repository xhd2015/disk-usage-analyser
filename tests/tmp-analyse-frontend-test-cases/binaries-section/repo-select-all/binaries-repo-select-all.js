const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

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

const firstRepoBlock = await page.$('[data-testid="binaries-tree"] > div');
if (!firstRepoBlock) {
    console.log('SKIP repo-select-all: no repo rows');
    console.log('DONE');
    process.exit(0);
}

const repoCheckbox = await firstRepoBlock.$('[data-testid^="binary-repo-checkbox-"]');
if (!repoCheckbox) {
    console.log('SKIP repo-select-all: no repo rows');
    console.log('DONE');
    process.exit(0);
}

const leafCount = await firstRepoBlock.$$eval('[data-testid^="binary-checkbox-"]', els => els.length);
console.log(`COUNT binary-leaf-checkboxes: ${leafCount}`);

if (leafCount < 2) {
    console.log('SKIP repo-select-all: need at least 2 binaries in a repo');
    console.log('DONE');
    process.exit(0);
}

await repoCheckbox.click();
await page.waitForTimeout(300);

const checked = await firstRepoBlock.$$eval('[data-testid^="binary-checkbox-"]', els =>
    els.filter(el => el.checked || el.getAttribute('aria-checked') === 'true').length
);
console.log(`COUNT binary-checkboxes-checked: ${checked}`);

const total = await page.$eval('[data-testid="binary-selected-total"]', el => el.textContent.trim()).catch(() => '');
console.log(`SELECTED_TOTAL: "${total}"`);

const allSelected = checked === leafCount;
console.log(`CHECK repo-select-all: ${allSelected}`);

if (!allSelected) process.exitCode = 1;
console.log('DONE');