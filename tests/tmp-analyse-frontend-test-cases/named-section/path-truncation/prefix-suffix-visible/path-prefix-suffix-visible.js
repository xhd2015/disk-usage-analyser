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

const pathEls = await page.$$('[data-testid^="node-modules-path-"]');
console.log(`ELEM node-modules-path: ${pathEls.length > 0 ? 'present' : 'MISSING'}`);

function startsWithEllipsis(text) {
    return text.startsWith('…') || text.startsWith('...');
}

function endsWithNodeModules(text) {
    return text.endsWith('node_modules') || text.endsWith('/node_modules');
}

if (namedRows.length > 0 && pathEls.length > 0) {
    const paths = await page.$$eval('[data-testid^="node-modules-path-"]', els =>
        els.map(el => ({
            visible: (el.textContent || '').trim(),
            full: (el.getAttribute('data-full-path') || '').trim(),
        }))
    ).catch(() => []);

    const longPath = paths.find(p => p.full.length > 40 && p.full.includes('node_modules'))
        || paths.find(p => p.visible.length > 40 && p.visible.includes('node_modules'));

    if (!longPath) {
        console.log('PATH_LONG: none');
    } else {
        const visible = longPath.visible;
        const full = longPath.full || visible;
        console.log(`PATH_VISIBLE_TEXT: "${visible}"`);
        console.log(`PATH_FULL_ATTR_TEXT: "${full}"`);

        const prefixOk = startsWithEllipsis(visible);
        const suffixOk = endsWithNodeModules(visible);
        const fullAttrOk = full.length > 0
            && full.includes('node_modules')
            && full.length >= visible.replace(/^[….]{1,3}/, '').length;

        console.log(`PATH_PREFIX_ELLIPSIS: ${prefixOk ? 'ok' : 'fail'}`);
        console.log(`PATH_SUFFIX_NODE_MODULES: ${suffixOk ? 'ok' : 'fail'}`);
        console.log(`PATH_FULL_ATTR: ${fullAttrOk ? 'ok' : 'fail'}`);

        if (!prefixOk || !suffixOk || !fullAttrOk) {
            process.exitCode = 1;
        }
    }
}

const empty = await page.$('[data-testid="node-modules-empty-state"]');
if (namedRows.length === 0 && empty) {
    console.log('NAMED_EMPTY_STATE: present');
}

if (namedRows.length > 0 && pathEls.length === 0) process.exitCode = 1;
console.log('DONE');