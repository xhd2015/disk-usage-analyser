const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.setViewportSize({ width: 1280, height: 800 });
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

const header = await page.$('[data-testid="node-modules-column-header"]');
console.log(`ELEM node-modules-column-header: ${header ? 'present' : 'MISSING'}`);

const namedRows = await page.$$('[data-testid^="node-modules-row-"]');
console.log(`COUNT node-modules-rows: ${namedRows.length}`);

if (tree && header && namedRows.length > 0) {
    const metrics = await page.evaluate(() => {
        const section = document.querySelector('[data-testid="node-modules-section"]');
        const body = section?.querySelector('.ant-card-body');
        const headerEl = document.querySelector('[data-testid="node-modules-column-header"]');
        const sizeCell = document.querySelector('[data-testid^="node-modules-size-"]');
        if (!body || !headerEl) {
            return null;
        }
        const bodyRect = body.getBoundingClientRect();
        const headerRect = headerEl.getBoundingClientRect();
        const sizeRect = sizeCell?.getBoundingClientRect();
        const headerRatio = headerRect.width / bodyRect.width;
        const sizeRightRatio = sizeRect
            ? (sizeRect.right - bodyRect.left) / bodyRect.width
            : headerRatio;
        return {
            bodyWidth: Math.round(bodyRect.width),
            headerWidth: Math.round(headerRect.width),
            sizeRight: sizeRect ? Math.round(sizeRect.right - bodyRect.left) : 0,
            headerRatio,
            sizeRightRatio,
        };
    });

    if (!metrics) {
        console.log('WIDTH_METRICS: MISSING');
        process.exitCode = 1;
    } else {
        console.log(`WIDTH body: ${metrics.bodyWidth}px`);
        console.log(`WIDTH header: ${metrics.headerWidth}px`);
        console.log(`WIDTH size-right-offset: ${metrics.sizeRight}px`);
        console.log(`WIDTH_RATIO header/body: ${metrics.headerRatio.toFixed(3)}`);
        console.log(`WIDTH_RATIO size-right/body: ${metrics.sizeRightRatio.toFixed(3)}`);
        if (metrics.sizeRightRatio < 0.90) {
            process.exitCode = 1;
        }
    }
} else if (!tree) {
    process.exitCode = 1;
} else if (namedRows.length === 0) {
    const empty = await page.$('[data-testid="node-modules-empty-state"]');
    if (empty) {
        console.log('NAMED_EMPTY_STATE: present');
    }
}

console.log('DONE');