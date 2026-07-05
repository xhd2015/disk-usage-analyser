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

const tree = await page.$('[data-testid="node-modules-tree"]');
console.log(`ELEM node-modules-tree: ${tree ? 'present' : 'MISSING'}`);

const namedRows = await page.$$('[data-testid^="node-modules-row-"]');
console.log(`COUNT node-modules-rows: ${namedRows.length}`);

const header = await page.$('[data-testid="node-modules-column-header"]');
console.log(`ELEM node-modules-column-header: ${header ? 'present' : 'MISSING'}`);

const pkgMgr = await page.$('[data-testid^="node-modules-pkgmgr-"]');
console.log(`ELEM node-modules-pkgmgr: ${pkgMgr ? 'present' : 'MISSING'}`);

const shared = await page.$('[data-testid^="node-modules-shared-"]');
console.log(`ELEM node-modules-shared: ${shared ? 'present' : 'MISSING'}`);

if (namedRows.length > 0) {
    const pkgMgrValue = await page.$eval('[data-testid^="node-modules-pkgmgr-"]', el => el.textContent.trim()).catch(() => '');
    console.log(`PKG_MGR_VALUE: "${pkgMgrValue}"`);
    const sharedValue = await page.$eval('[data-testid^="node-modules-shared-"]', el => el.textContent.trim()).catch(() => '');
    console.log(`SHARED_VALUE: "${sharedValue}"`);
}

const empty = await page.$('[data-testid="node-modules-empty-state"]');
if (namedRows.length === 0 && empty) {
    console.log('NAMED_EMPTY_STATE: present');
}

if (!tree) process.exitCode = 1;
if (namedRows.length > 0) {
    if (!header || !pkgMgr || !shared) process.exitCode = 1;
}
console.log('DONE');