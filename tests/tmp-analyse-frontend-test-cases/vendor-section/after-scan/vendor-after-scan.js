const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

const scanBtn = await page.$('[data-testid="vendor-scan-btn"]');
if (!scanBtn) {
    console.log('FAIL vendor-scan-btn: MISSING');
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

await scanBtn.click();

for (let i = 0; i < 120; i++) {
    await page.waitForTimeout(500);
    const done = await page.$('[data-testid="vendor-section"] [data-testid="vendor-done-badge"]');
    if (done) break;
}

const tree = await page.$('[data-testid="vendor-tree"]');
console.log(`ELEM vendor-tree: ${tree ? 'present' : 'MISSING'}`);

const vendorRows = await page.$$('[data-testid^="vendor-row-"]');
console.log(`COUNT vendor-rows: ${vendorRows.length}`);

if (vendorRows.length > 0) {
    const size = await page.$eval('[data-testid^="vendor-size-"]', el => el.textContent.trim()).catch(() => '');
    console.log(`VENDOR_ENTITY_SIZE: "${size}"`);
    const name = await page.$eval('[data-testid^="vendor-name-"]', el => el.textContent.trim()).catch(() => '');
    console.log(`VENDOR_ENTITY_NAME: "${name}"`);
    const path = await page.$eval('[data-testid^="vendor-path-"]', el => el.textContent.trim()).catch(() => '');
    console.log(`VENDOR_ENTITY_PATH: "${path}"`);
    const repo = await page.$eval('[data-testid^="vendor-repo-"]', el => el.textContent.trim()).catch(() => '');
    console.log(`VENDOR_ENTITY_REPO: "${repo}"`);
}

const empty = await page.$('[data-testid="vendor-empty-state"]');
if (vendorRows.length === 0 && empty) {
    console.log('VENDOR_EMPTY_STATE: present');
}

if (!tree) process.exitCode = 1;
console.log('DONE');
