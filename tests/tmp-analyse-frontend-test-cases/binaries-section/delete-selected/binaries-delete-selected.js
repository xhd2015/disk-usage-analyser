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

const rowsBefore = (await page.$$('[data-testid^="binary-row-"]')).length;
console.log(`COUNT binary-rows-before: ${rowsBefore}`);

if (rowsBefore === 0) {
    console.log('SKIP binaries-delete: no binary rows');
    console.log('DONE');
    process.exit(0);
}

const checkbox = await page.$('[data-testid^="binary-checkbox-"]');
const pathTestId = await checkbox.getAttribute('data-testid');
console.log(`SELECTED checkbox: ${pathTestId}`);

await checkbox.click();
await page.waitForTimeout(200);

const deleteBtn = await page.$('[data-testid="binary-delete-btn"]');
if (!deleteBtn) {
    console.log('FAIL binary-delete-btn: MISSING');
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}
await deleteBtn.click();
await page.waitForTimeout(300);

const modal = await page.$('[data-testid="binary-delete-confirm-modal"]');
console.log(`ELEM binary-delete-confirm-modal: ${modal ? 'present' : 'MISSING'}`);

const confirmBtn = await page.$('[data-testid="binary-delete-confirm-btn"]');
if (!confirmBtn) {
    console.log('FAIL binary-delete-confirm-btn: MISSING');
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

await confirmBtn.click();

for (let i = 0; i < 30; i++) {
    await page.waitForTimeout(500);
    const success = await page.$('[data-testid="binary-delete-success"]');
    if (success) break;
}

const rowsAfter = (await page.$$('[data-testid^="binary-row-"]')).length;
console.log(`COUNT binary-rows-after: ${rowsAfter}`);
console.log(`CHECK binary-row-removed: ${rowsAfter < rowsBefore}`);

if (!modal || rowsAfter >= rowsBefore) process.exitCode = 1;
console.log('DONE');