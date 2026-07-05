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

const namedRows = await page.$$('[data-testid^="node-modules-row-"]');
console.log(`COUNT node-modules-rows: ${namedRows.length}`);

const pathEl = await page.$('[data-testid^="node-modules-path-"]');
console.log(`ELEM node-modules-path: ${pathEl ? 'present' : 'MISSING'}`);

if (namedRows.length > 0 && pathEl) {
    const pathInfo = await page.$eval('[data-testid^="node-modules-path-"]', el => {
        const style = window.getComputedStyle(el);
        const visible = (el.textContent || '').trim();
        const full = (el.getAttribute('data-full-path') || visible).trim();
        const hasTruncationStyle =
            style.textOverflow === 'ellipsis' ||
            style.overflow === 'hidden' ||
            visible.startsWith('…') ||
            visible.startsWith('...') ||
            el.closest('[style*="ellipsis"]') !== null;
        return {
            visible,
            full,
            hasTruncationStyle,
            overflow: style.overflow,
            textOverflow: style.textOverflow,
        };
    }).catch(() => ({ visible: '', full: '', hasTruncationStyle: false, overflow: '', textOverflow: '' }));

    console.log(`PATH_VISIBLE_TEXT: "${pathInfo.visible}"`);
    console.log(`PATH_FULL_ATTR_TEXT: "${pathInfo.full}"`);
    console.log(`PATH_ELLIPSIS: ${pathInfo.hasTruncationStyle ? 'ok' : 'fail'}`);
    console.log(`PATH_FULL_ATTR: ${pathInfo.full.length > 20 ? 'ok' : 'fail'}`);

    await pathEl.hover();
    await page.waitForTimeout(400);

    const tooltipText = await page.$eval('.ant-tooltip-inner', el => el.textContent.trim()).catch(() => '');
    const tooltipOk = tooltipText.length > 0
        && pathInfo.full.length > 0
        && tooltipText.includes(pathInfo.full);
    console.log(`PATH_TOOLTIP: ${tooltipOk ? 'ok' : 'fail'}`);
    if (!tooltipOk) console.log(`PATH_TOOLTIP_TEXT: "${tooltipText}"`);

    if (!pathInfo.hasTruncationStyle || pathInfo.full.length <= 20 || !tooltipOk) {
        process.exitCode = 1;
    }
}

const empty = await page.$('[data-testid="node-modules-empty-state"]');
if (namedRows.length === 0 && empty) {
    console.log('NAMED_EMPTY_STATE: present');
}

if (namedRows.length > 0 && !pathEl) process.exitCode = 1;
console.log('DONE');