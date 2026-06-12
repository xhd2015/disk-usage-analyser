const BASE = process.env.SERVER_URL;

// 1. Verify server is reachable via /ping
const pingResp = await page.goto(`${BASE}/ping`, { waitUntil: 'load' });
const pingBody = await page.evaluate(() => document.body.textContent);
console.log(`PING: "${pingBody}"`);
console.log(`PING_OK: ${pingBody === 'pong'}`);

// 2. Navigate to tmp-analyse and verify basic structure renders
await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

const heading = await page.$('[data-testid="page-heading"]');
console.log(`PAGE_HEADING present: ${heading !== null}`);

const startBtn = await page.$('[data-testid="start-scan-btn"]');
console.log(`START_BTN present: ${startBtn !== null}`);

// 3. Check that core cards exist
for (const cat of ['trash', 'temp', 'cache', 'log']) {
    const card = await page.$(`[data-testid="card-${cat}"]`);
    console.log(`CARD ${cat} present: ${card !== null}`);
}

// 4. Verify SERVER_URL is set (not defaulted to localhost:8080)
console.log(`SERVER_URL: ${BASE}`);

console.log('DONE');
