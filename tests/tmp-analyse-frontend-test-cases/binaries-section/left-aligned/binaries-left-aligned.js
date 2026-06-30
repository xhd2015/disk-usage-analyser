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
if (!tree) {
    console.log('ELEM binaries-tree: MISSING');
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

console.log('ELEM binaries-tree: present');

const textAlign = await tree.evaluate(el => getComputedStyle(el).textAlign);
const width = await tree.evaluate(el => getComputedStyle(el).width);
console.log(`TREE_STYLE text-align: ${textAlign}`);
console.log(`TREE_STYLE width: ${width}`);

if (textAlign !== 'left') process.exitCode = 1;
console.log('DONE');