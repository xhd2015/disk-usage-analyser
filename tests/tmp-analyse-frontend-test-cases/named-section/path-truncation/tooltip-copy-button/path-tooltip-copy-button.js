const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

try {
    if (typeof context !== 'undefined' && context.grantPermissions) {
        await context.grantPermissions(['clipboard-read', 'clipboard-write']);
    } else if (page.context && page.context().grantPermissions) {
        await page.context().grantPermissions(['clipboard-read', 'clipboard-write']);
    }
} catch (_) {
    // clipboard permissions may already be granted in playwright-debug
}

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
    const fullPath = await page.$eval('[data-testid^="node-modules-path-"]', el =>
        (el.getAttribute('data-full-path') || el.textContent || '').trim()
    ).catch(() => '');
    console.log(`PATH_FULL_ATTR_TEXT: "${fullPath}"`);

    await pathEl.hover();
    await page.waitForTimeout(500);

    const copyBtn = await page.$('[data-testid^="node-modules-path-copy-"]');
    console.log(`ELEM node-modules-path-copy: ${copyBtn ? 'present' : 'MISSING'}`);

    let clipboardText = '';
    if (copyBtn && fullPath) {
        await copyBtn.click();
        await page.waitForTimeout(300);
        clipboardText = await page.evaluate(async () => {
            try {
                return await navigator.clipboard.readText();
            } catch (e) {
                return '';
            }
        }).catch(() => '');
    }

    console.log(`PATH_CLIPBOARD_TEXT: "${clipboardText}"`);
    const clipboardOk = clipboardText.length > 0 && clipboardText === fullPath;
    console.log(`PATH_CLIPBOARD: ${clipboardOk ? 'ok' : 'fail'}`);

    if (!copyBtn || !clipboardOk) {
        process.exitCode = 1;
    }
}

const empty = await page.$('[data-testid="node-modules-empty-state"]');
if (namedRows.length === 0 && empty) {
    console.log('NAMED_EMPTY_STATE: present');
}

if (namedRows.length > 0 && !pathEl) process.exitCode = 1;
console.log('DONE');