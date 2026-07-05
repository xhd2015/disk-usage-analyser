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

let rowsAtScanComplete = false;
let sharedAtScanComplete = '';
let sharedFinal = '';
let maxRows = 0;

await scanBtn.click();

for (let i = 0; i < 120; i++) {
    await page.waitForTimeout(500);

    const namedRows = await page.$$('[data-testid^="node-modules-row-"]');
    if (namedRows.length > maxRows) {
        maxRows = namedRows.length;
    }

    const doneBadge = await page.$('[data-testid="node-modules-section"] [data-testid="node-modules-done-badge"]');
    if (doneBadge && namedRows.length > 0 && !rowsAtScanComplete) {
        rowsAtScanComplete = true;
        sharedAtScanComplete = await page.$eval('[data-testid^="node-modules-shared-"]', el => el.textContent.trim()).catch(() => '');
        console.log(`TIMING rows-at-scan-complete: count=${namedRows.length}`);
        console.log(`SHARED_AT_SCAN_COMPLETE: "${sharedAtScanComplete}"`);
    }

    const enriching = await page.$('[data-testid="node-modules-enriching-badge"]');
    if (doneBadge && !enriching) {
        sharedFinal = await page.$eval('[data-testid^="node-modules-shared-"]', el => el.textContent.trim()).catch(() => '');
        break;
    }
}

if (!sharedFinal && maxRows > 0) {
    sharedFinal = await page.$eval('[data-testid^="node-modules-shared-"]', el => el.textContent.trim()).catch(() => '');
}

const namedRows = await page.$$('[data-testid^="node-modules-row-"]');
console.log(`COUNT node-modules-rows: ${namedRows.length}`);

const empty = await page.$('[data-testid="node-modules-empty-state"]');
if (namedRows.length === 0 && empty) {
    console.log('NAMED_EMPTY_STATE: present');
}

console.log(`CHECK rows-at-scan-complete: ${rowsAtScanComplete}`);
console.log(`SHARED_FINAL: "${sharedFinal}"`);
const sharedFinalPresent = sharedFinal.length > 0;
console.log(`CHECK shared-final-present: ${sharedFinalPresent}`);

if (namedRows.length > 0) {
    if (!rowsAtScanComplete || !sharedFinalPresent) process.exitCode = 1;
}

console.log('DONE');