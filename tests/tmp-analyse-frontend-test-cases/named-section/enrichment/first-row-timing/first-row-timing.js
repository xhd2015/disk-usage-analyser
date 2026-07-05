const BASE = process.env.SERVER_URL || 'http://localhost:8080';

const t0 = Date.now();

await page.addInitScript(() => {
    window.__nmTimings = { firstNamedMs: null, scanCompleteMs: null, doneMs: null };
    const OrigES = window.EventSource;
    window.EventSource = function(url, opts) {
        const es = new OrigES(url, opts);
        if (url.includes('tmp-named') && url.includes('node_modules')) {
            es.addEventListener('named', () => {
                if (window.__nmTimings.firstNamedMs === null) {
                    window.__nmTimings.firstNamedMs = performance.now();
                }
            });
            es.addEventListener('scan_complete', () => {
                window.__nmTimings.scanCompleteMs = performance.now();
            });
            es.addEventListener('done', () => {
                window.__nmTimings.doneMs = performance.now();
            });
        }
        return es;
    };
});

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

const scanBtn = await page.$('[data-testid="node-modules-scan-btn"]');
if (!scanBtn) {
    console.log('FAIL node-modules-scan-btn: MISSING');
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

let firstRowMs = null;

await scanBtn.click();
const clickMs = Date.now();

for (let i = 0; i < 240; i++) {
    await page.waitForTimeout(500);

    const rows = await page.$$('[data-testid^="node-modules-row-"]');
    if (rows.length > 0 && firstRowMs === null) {
        firstRowMs = Date.now() - clickMs;
        console.log(`TIMING first-row-ui-ms: ${firstRowMs}`);
    }

    const timings = await page.evaluate(() => window.__nmTimings || {});
    if (timings.doneMs != null) break;
}

const timings = await page.evaluate(() => window.__nmTimings || {});
const namedRows = await page.$$('[data-testid^="node-modules-row-"]');

console.log(`COUNT node-modules-rows: ${namedRows.length}`);
if (timings.firstNamedMs != null) {
    console.log(`TIMING first-named-sse-ms: ${Math.round(timings.firstNamedMs)}`);
}
if (timings.scanCompleteMs != null) {
    console.log(`TIMING scan-complete-sse-ms: ${Math.round(timings.scanCompleteMs)}`);
}
if (timings.doneMs != null) {
    console.log(`TIMING done-sse-ms: ${Math.round(timings.doneMs)}`);
    if (timings.scanCompleteMs != null) {
        console.log(`TIMING enrichment-phase-ms: ${Math.round(timings.doneMs - timings.scanCompleteMs)}`);
    }
}

const FIRST_ROW_BUDGET_MS = 10000;
if (firstRowMs != null) {
    console.log(`CHECK first-row-within-budget: ${firstRowMs <= FIRST_ROW_BUDGET_MS}`);
    if (firstRowMs > FIRST_ROW_BUDGET_MS) process.exitCode = 1;
} else if (namedRows.length === 0) {
    const empty = await page.$('[data-testid="node-modules-empty-state"]');
    if (empty) console.log('NAMED_EMPTY_STATE: present');
    else process.exitCode = 1;
} else {
    process.exitCode = 1;
}

console.log('DONE');