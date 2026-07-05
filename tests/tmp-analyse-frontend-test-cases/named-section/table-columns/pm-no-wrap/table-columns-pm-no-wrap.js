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

const pmCells = await page.$$('[data-testid^="node-modules-pkgmgr-"]');
console.log(`COUNT node-modules-pkgmgr: ${pmCells.length}`);
console.log(`ELEM node-modules-pkgmgr: ${pmCells.length > 0 ? 'present' : 'MISSING'}`);

if (namedRows.length > 0 && pmCells.length > 0) {
    const styles = await page.$$eval('[data-testid^="node-modules-pkgmgr-"]', els =>
        els.map(el => window.getComputedStyle(el).whiteSpace)
    ).catch(() => []);

    const allNowrap = styles.length > 0 && styles.every(s => s === 'nowrap');
    console.log(`PM_NOWRAP: ${allNowrap ? 'ok' : 'fail'}`);
    console.log(`PM_WHITE_SPACE_VALUES: ${JSON.stringify(styles)}`);

    const pmValue = await page.$eval('[data-testid^="node-modules-pkgmgr-"]', el => (el.textContent || '').trim()).catch(() => '');
    console.log(`PKG_MGR_VALUE: "${pmValue}"`);

    if (!allNowrap) process.exitCode = 1;
}

const empty = await page.$('[data-testid="node-modules-empty-state"]');
if (namedRows.length === 0 && empty) {
    console.log('NAMED_EMPTY_STATE: present');
}

if (!tree) process.exitCode = 1;
if (namedRows.length > 0) {
    if (!header || pmCells.length === 0) process.exitCode = 1;
}
console.log('DONE');