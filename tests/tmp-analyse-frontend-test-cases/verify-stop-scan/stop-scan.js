const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

// Click Start Scan
const startBtn = await page.$('[data-testid="start-scan-btn"]');
if (!startBtn) { console.log('BUTTON start-scan-btn: MISSING'); process.exitCode = 1; }
await startBtn.click();
await page.waitForTimeout(1000);

// Verify scan started (stop button visible)
try {
    const sv = await page.$eval('[data-testid="stop-scan-btn"]', el => {
        return window.getComputedStyle(el).display !== 'none';
    });
    console.log(`BUTTON scan-started stop-scan-btn: visible=${sv}`);
    if (!sv) {
        console.log('FAIL: stop-scan-btn not visible after starting scan');
        process.exitCode = 1;
    }
} catch {
    console.log('BUTTON scan-started stop-scan-btn: visible=false');
    console.log('FAIL: stop-scan-btn not visible after starting scan');
    process.exitCode = 1;
}

// Click Stop Scan
const stopBtn = await page.$('[data-testid="stop-scan-btn"]');
if (stopBtn) {
    await stopBtn.click();
    await page.waitForTimeout(500);

    // Verify scan stopped
    try {
        const sv = await page.$eval('[data-testid="start-scan-btn"]', el => {
            return window.getComputedStyle(el).display !== 'none';
        });
        console.log(`BUTTON after stop start-scan-btn: visible=${sv}`);
    } catch {
        console.log('BUTTON after stop start-scan-btn: visible=false');
    }
    try {
        const sv = await page.$eval('[data-testid="stop-scan-btn"]', el => {
            return window.getComputedStyle(el).display !== 'none';
        });
        console.log(`BUTTON after stop stop-scan-btn: visible=${sv}`);
    } catch {
        console.log('BUTTON after stop stop-scan-btn: visible=false');
    }
} else {
    console.log('FAIL: stop-scan-btn not found to click');
    process.exitCode = 1;
}

console.log('DONE');
