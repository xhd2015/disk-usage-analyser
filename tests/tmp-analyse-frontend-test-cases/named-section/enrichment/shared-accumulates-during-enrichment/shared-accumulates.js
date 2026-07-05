const BASE = process.env.SERVER_URL || 'http://localhost:8080';

function parseSharedBytes(text) {
    const t = text.trim();
    if (!t || t === '…' || t.includes('loading')) return 0;
    const m = t.match(/^([\d.]+)\s*([KMGT]?)/);
    if (!m) return 0;
    const mult = { '': 1, K: 1024, M: 1024 ** 2, G: 1024 ** 3, T: 1024 ** 4 }[m[2]] || 1;
    return parseFloat(m[1]) * mult;
}

async function sharedSum() {
    const texts = await page.$$eval('[data-testid^="node-modules-shared-"]', els =>
        els.map(el => el.textContent.trim())
    ).catch(() => []);
    return texts.reduce((sum, t) => sum + parseSharedBytes(t), 0);
}

await page.addInitScript(() => {
    window.__nmEnrichedCount = 0;
    const OrigES = window.EventSource;
    window.EventSource = function(url, opts) {
        const es = new OrigES(url, opts);
        if (url.includes('tmp-named') && url.includes('node_modules')) {
            es.addEventListener('named_enriched', () => {
                window.__nmEnrichedCount++;
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

const t0 = Date.now();
let rowUpdates = 0;
let prevSharedTexts = [];
let maxEnrichedSSE = 0;

await scanBtn.click();

for (let i = 0; i < 240; i++) {
    await page.waitForTimeout(500);

    const enrichedSSE = await page.evaluate(() => window.__nmEnrichedCount || 0);
    if (enrichedSSE > maxEnrichedSSE) maxEnrichedSSE = enrichedSSE;

    const rows = await page.$$('[data-testid^="node-modules-row-"]');
    if (rows.length > 0 && enrichedSSE > 0) {
        const sharedTexts = await page.$$eval('[data-testid^="node-modules-shared-"]', els =>
            els.map(el => el.textContent.trim())
        ).catch(() => []);
        let changed = 0;
        for (let j = 0; j < sharedTexts.length; j++) {
            if (prevSharedTexts[j] !== undefined && sharedTexts[j] !== prevSharedTexts[j]) {
                changed++;
            }
        }
        if (changed > 0) {
            rowUpdates++;
            const sum = await sharedSum();
            console.log(`SHARED_ROW_UPDATES #${rowUpdates}: changed=${changed} sum=${Math.round(sum)} enriched=${enrichedSSE} at ${Date.now() - t0}ms`);
        }
        prevSharedTexts = sharedTexts;
    }

    const doneBadge = await page.$('[data-testid="node-modules-done-badge"]');
    const enriching = await page.$('[data-testid="node-modules-enriching-badge"]');
    if (doneBadge && !enriching && rows.length > 0 && enrichedSSE >= rows.length) break;
}

const namedRows = await page.$$('[data-testid^="node-modules-row-"]');
console.log(`COUNT node-modules-rows: ${namedRows.length}`);
console.log(`METRIC enriched-sse-events: ${maxEnrichedSSE}`);
console.log(`METRIC shared-row-updates: ${rowUpdates}`);
console.log(`CHECK shared-accumulating: ${rowUpdates >= 2}`);

const empty = await page.$('[data-testid="node-modules-empty-state"]');
if (namedRows.length === 0 && empty) {
    console.log('NAMED_EMPTY_STATE: present');
}

if (namedRows.length > 0 && maxEnrichedSSE >= 2) {
    if (rowUpdates < 2) process.exitCode = 1;
}

console.log('DONE');