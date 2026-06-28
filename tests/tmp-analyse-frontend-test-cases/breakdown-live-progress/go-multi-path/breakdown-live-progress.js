const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

const goCard = await page.$('[data-testid="card-go"]');
if (!goCard) {
    console.log('SKIP breakdown-live: Go card not detected');
    console.log('DONE');
    process.exit(0);
}

const startBtn = await page.$('[data-testid="start-scan-btn"]');
if (!startBtn) {
    console.log('FAIL start-scan-btn: MISSING');
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

await startBtn.click();

let sawBreakdown0NonZero = false;
let sawBreakdown1Update = false;
let sawScanningBeforeDone = false;

for (let i = 0; i < 30; i++) {
    await page.waitForTimeout(200);

    const scanning = await page.$('[data-testid="card-go"] [data-testid="scanning-badge"]');
    const done = await page.$('[data-testid="card-go"] [data-testid="done-badge"]');

    const size0 = await page.$eval(
        '[data-testid="card-go"] [data-testid="breakdown-size-0"]',
        el => (el.textContent || '').trim()
    ).catch(() => '');

    const size1 = await page.$eval(
        '[data-testid="card-go"] [data-testid="breakdown-size-1"]',
        el => (el.textContent || '').trim()
    ).catch(() => '');

    if (scanning) {
        sawScanningBeforeDone = true;
        if (size0 && size0 !== '0 Bytes' && size0 !== '0 B' && size0 !== '-') {
            sawBreakdown0NonZero = true;
            console.log(`LIVE breakdown-size-0: "${size0}"`);
        }
        if (size1 && size1 !== '0 Bytes' && size1 !== '0 B' && size1 !== '-') {
            sawBreakdown1Update = true;
            console.log(`LIVE breakdown-size-1: "${size1}"`);
        }
    }

    if (done) {
        console.log(`FINAL breakdown-size-0: "${size0}"`);
        console.log(`FINAL breakdown-size-1: "${size1}"`);
        break;
    }
}

console.log(`CHECK scanning-before-done: ${sawScanningBeforeDone}`);
console.log(`CHECK breakdown-0-live-update: ${sawBreakdown0NonZero}`);
console.log(`CHECK breakdown-1-live-update: ${sawBreakdown1Update}`);

if (!sawScanningBeforeDone) {
    console.log('FAIL: never saw scanning badge on Go card');
    process.exitCode = 1;
}
if (!sawBreakdown0NonZero && !sawBreakdown1Update) {
    console.log('FAIL: breakdown rows did not update during scan');
    process.exitCode = 1;
}

console.log('DONE');