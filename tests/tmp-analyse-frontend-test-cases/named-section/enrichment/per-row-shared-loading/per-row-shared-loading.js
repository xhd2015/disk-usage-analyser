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

let perRowLoadingSeen = false;
let pendingSharedRowsMax = 0;

await page.addInitScript(() => {
    window.__nmEnrichedPaths = new Set();
    const OrigES = window.EventSource;
    window.EventSource = function(url, opts) {
        const es = new OrigES(url, opts);
        if (url.includes('tmp-named') && url.includes('node_modules')) {
            es.addEventListener('named_enriched', (e) => {
                try {
                    const hit = JSON.parse(e.data);
                    if (hit.path) window.__nmEnrichedPaths.add(hit.path);
                } catch (_) {}
            });
        }
        return es;
    };
});

await scanBtn.click();

for (let i = 0; i < 180; i++) {
    await page.waitForTimeout(500);

    const rowLoaders = await page.$$('[data-testid^="node-modules-shared-loading-"]');
    if (rowLoaders.length > 0) {
        perRowLoadingSeen = true;
        console.log(`ROW_LOADING count=${rowLoaders.length}`);
    }

    const namedRows = await page.$$('[data-testid^="node-modules-row-"]');
    const enrichedCount = await page.evaluate(() => window.__nmEnrichedPaths?.size || 0);
    if (namedRows.length > enrichedCount) {
        const pending = namedRows.length - enrichedCount;
        if (pending > pendingSharedRowsMax) pendingSharedRowsMax = pending;
    }

    const doneBadge = await page.$('[data-testid="node-modules-done-badge"]');
    const enrichingBadge = await page.$('[data-testid="node-modules-enriching-badge"]');
    if (doneBadge && !enrichingBadge && namedRows.length > 0 && enrichedCount >= namedRows.length) {
        break;
    }
}

const namedRows = await page.$$('[data-testid^="node-modules-row-"]');
console.log(`COUNT node-modules-rows: ${namedRows.length}`);
console.log(`METRIC pending-shared-rows-max: ${pendingSharedRowsMax}`);
console.log(`CHECK per-row-shared-loading-seen: ${perRowLoadingSeen}`);

const empty = await page.$('[data-testid="node-modules-empty-state"]');
if (namedRows.length === 0 && empty) {
    console.log('NAMED_EMPTY_STATE: present');
}

if (namedRows.length > 0 && pendingSharedRowsMax > 0) {
    if (!perRowLoadingSeen) process.exitCode = 1;
}

console.log('DONE');