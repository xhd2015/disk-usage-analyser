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

const repoRows = await page.$$('[data-testid^="node-modules-repo-row-"]');
console.log(`COUNT node-modules-repo-rows: ${repoRows.length}`);

const namedRows = await page.$$('[data-testid^="node-modules-row-"]');
console.log(`COUNT node-modules-rows: ${namedRows.length}`);

if (namedRows.length > 0) {
    const size = await page.$eval('[data-testid^="node-modules-size-"]', el => el.textContent.trim()).catch(() => '');
    console.log(`NAMED_ENTITY_SIZE: "${size}"`);
    const name = await page.$eval('[data-testid^="node-modules-name-"]', el => el.textContent.trim()).catch(() => '');
    console.log(`NAMED_ENTITY_NAME: "${name}"`);
    const path = await page.$eval('[data-testid^="node-modules-path-"]', el => el.textContent.trim()).catch(() => '');
    console.log(`NAMED_ENTITY_PATH: "${path}"`);
    const repo = await page.$eval('[data-testid^="node-modules-repo-"]', el => el.textContent.trim()).catch(() => '');
    console.log(`NAMED_ENTITY_REPO: "${repo}"`);
}

const empty = await page.$('[data-testid="node-modules-empty-state"]');
if (namedRows.length === 0 && empty) {
    console.log('NAMED_EMPTY_STATE: present');
}

if (!tree) process.exitCode = 1;
console.log('DONE');
