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

if (tree) {
    const style = await tree.evaluate(el => window.getComputedStyle(el).textAlign);
    console.log(`TREE_STYLE text-align: ${style}`);
} else {
    process.exitCode = 1;
}

console.log('DONE');
