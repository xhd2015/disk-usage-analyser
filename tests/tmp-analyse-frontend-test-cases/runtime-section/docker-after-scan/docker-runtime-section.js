const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

const dockerCard = await page.$('[data-testid="card-docker"]');
if (!dockerCard) {
    console.log('SKIP docker-runtime: Docker card not detected');
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

// Wait for Docker scan to complete
for (let i = 0; i < 60; i++) {
    await page.waitForTimeout(500);
    const done = await page.$('[data-testid="card-docker"] [data-testid="done-badge"]');
    if (done) break;
}

const runtimeSection = await page.$('[data-testid="card-docker"] [data-testid="runtime-section"]');
if (runtimeSection) {
    console.log('DOCKER_RUNTIME_SECTION: present');
    const row0 = await page.$('[data-testid="card-docker"] [data-testid="runtime-row-0"]');
    const label0 = await page.$('[data-testid="card-docker"] [data-testid="runtime-label-0"]');
    const count0 = await page.$('[data-testid="card-docker"] [data-testid="runtime-count-0"]');
    const size0 = await page.$('[data-testid="card-docker"] [data-testid="runtime-size-0"]');
    console.log(`ELEM runtime-row-0: ${row0 ? 'present' : 'MISSING'}`);
    console.log(`ELEM runtime-label-0: ${label0 ? (await label0.textContent()).trim() : 'MISSING'}`);
    console.log(`ELEM runtime-count-0: ${count0 ? (await count0.textContent()).trim() : 'MISSING'}`);
    console.log(`ELEM runtime-size-0: ${size0 ? (await size0.textContent()).trim() : 'MISSING'}`);
} else {
    console.log('DOCKER_RUNTIME_SECTION: absent');
    console.log('DOCKER_RUNTIME_GRACEFUL: true');
}

console.log('DONE');