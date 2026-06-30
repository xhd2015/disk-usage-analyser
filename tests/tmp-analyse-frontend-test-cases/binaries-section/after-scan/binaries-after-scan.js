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

const tree = await page.$('[data-testid="binaries-tree"]');
console.log(`ELEM binaries-tree: ${tree ? 'present' : 'MISSING'}`);

const repoRows = await page.$$('[data-testid^="binary-repo-row-"]');
console.log(`COUNT binary-repo-rows: ${repoRows.length}`);

const binRows = await page.$$('[data-testid^="binary-row-"]');
console.log(`COUNT binary-rows: ${binRows.length}`);

if (binRows.length > 0) {
    const kind = await page.$eval('[data-testid^="binary-kind-"]', el => el.textContent.trim()).catch(() => '');
    console.log(`BINARY_KIND: "${kind}"`);
    const path = await page.$eval('[data-testid^="binary-path-"]', el => el.textContent.trim()).catch(() => '');
    console.log(`BINARY_PATH: "${path}"`);
}

const empty = await page.$('[data-testid="binaries-empty-state"]');
if (binRows.length === 0 && empty) {
    console.log('BINARIES_EMPTY_STATE: present');
}

if (!tree) process.exitCode = 1;
console.log('DONE');