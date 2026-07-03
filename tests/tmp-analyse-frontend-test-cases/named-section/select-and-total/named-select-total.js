const BASE = process.env.SERVER_URL || 'http://localhost:8080';

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

const checkbox = await page.$('[data-testid^="node-modules-checkbox-"]');
if (!checkbox) {
    console.log('SKIP named-select: no node_modules rows');
    console.log('DONE');
    process.exit(0);
}

const totalBefore = await page.$eval('[data-testid="node-modules-selected-total"]', el => el.textContent.trim()).catch(() => '');
console.log(`SELECTED_TOTAL before: "${totalBefore}"`);

await checkbox.click();
await page.waitForTimeout(300);

const totalAfter = await page.$eval('[data-testid="node-modules-selected-total"]', el => el.textContent.trim()).catch(() => '');
console.log(`SELECTED_TOTAL after: "${totalAfter}"`);

const deleteBtn = await page.$('[data-testid="node-modules-delete-btn"]');
console.log(`ELEM node-modules-delete-btn: ${deleteBtn ? (await deleteBtn.isVisible() ? 'visible' : 'hidden') : 'MISSING'}`);

const changed = totalBefore !== totalAfter && totalAfter.includes('to clear');
console.log(`CHECK selected-total-updated: ${changed}`);

if (!changed) process.exitCode = 1;
console.log('DONE');
