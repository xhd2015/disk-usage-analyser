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

const rowsBefore = (await page.$$('[data-testid^="node-modules-row-"]')).length;
console.log(`COUNT named-rows-before: ${rowsBefore}`);

if (rowsBefore === 0) {
    console.log('SKIP named-delete: no node_modules rows');
    console.log('DONE');
    process.exit(0);
}

const checkbox = await page.$('[data-testid^="node-modules-checkbox-"]');
const pathTestId = await checkbox.getAttribute('data-testid');
console.log(`SELECTED checkbox: ${pathTestId}`);

await checkbox.click();
await page.waitForTimeout(200);

const deleteBtn = await page.$('[data-testid="node-modules-delete-btn"]');
if (!deleteBtn) {
    console.log('FAIL node-modules-delete-btn: MISSING');
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}
await deleteBtn.click();
await page.waitForTimeout(300);

const modal = await page.$('[data-testid="node-modules-delete-confirm-modal"]');
console.log(`ELEM node-modules-delete-confirm-modal: ${modal ? 'present' : 'MISSING'}`);

const confirmBtn = await page.$('[data-testid="node-modules-delete-confirm-btn"]');
if (!confirmBtn) {
    console.log('FAIL node-modules-delete-confirm-btn: MISSING');
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

await confirmBtn.click();

for (let i = 0; i < 30; i++) {
    await page.waitForTimeout(500);
    const success = await page.$('[data-testid="node-modules-delete-success"]');
    if (success) break;
}

const rowsAfter = (await page.$$('[data-testid^="node-modules-row-"]')).length;
console.log(`COUNT named-rows-after: ${rowsAfter}`);
console.log(`CHECK named-row-removed: ${rowsAfter < rowsBefore}`);

if (!modal || rowsAfter >= rowsBefore) process.exitCode = 1;
console.log('DONE');
