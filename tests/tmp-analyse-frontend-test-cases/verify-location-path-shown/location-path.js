const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

// Click Start Scan
const startBtn = await page.$('[data-testid="start-scan-btn"]');
if (!startBtn) { console.log('BUTTON start-scan-btn: MISSING'); process.exitCode = 1; }
await startBtn.click();

// Wait for scan to complete (start button reappears)
try {
    await page.waitForSelector('[data-testid="start-scan-btn"]', { state: 'visible', timeout: 15000 });
    console.log('SCAN_COMPLETE: start button reappeared');
} catch {
    console.log('SCAN_COMPLETE: timeout waiting for start button');
}

await page.waitForTimeout(500);

// Check that card-path elements exist and have non-empty text
const pathEls = await page.$$('[data-testid="card-path"]');
console.log(`PATH_ELEMENTS: ${pathEls.length}`);

const pathTexts = await page.$$eval('[data-testid="card-path"]', els => els.map(el => el.textContent));
console.log(`PATH_TEXTS: ${JSON.stringify(pathTexts)}`);

const hasNonEmpty = pathTexts.some(t => t && t.trim() !== '');
console.log(`HAS_NONEMPTY_PATH: ${hasNonEmpty}`);

const hasRecognizable = pathTexts.some(t => t && (t.includes('Trash') || t.includes('Caches') || t.includes('Logs') || t.includes('/tmp') || t.includes('Temp')));
console.log(`HAS_RECOGNIZABLE_PATH: ${hasRecognizable}`);

// Verify tilde prefix: home-relative paths should use ~ not /Users/
const homeRelativePaths = pathTexts.filter(t => t && !t.startsWith('/tmp') && !t.startsWith('/usr/'));
const hasTildePath = homeRelativePaths.some(t => t && t.startsWith('~/'));
console.log(`HAS_TILDE_PATH: ${hasTildePath}`);

const noRawHomeDir = !pathTexts.some(t => t && (t.includes('/Users/') || t.includes('/home/')));
console.log(`NO_RAW_HOME_DIR: ${noRawHomeDir}`);

console.log('DONE');
