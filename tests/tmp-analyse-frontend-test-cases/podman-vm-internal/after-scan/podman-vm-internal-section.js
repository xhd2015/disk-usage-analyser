const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

const podmanCard = await page.$('[data-testid="card-podman"]');
if (!podmanCard) {
    console.log('SKIP podman-vm-internal: Podman card not detected');
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

for (let i = 0; i < 60; i++) {
    await page.waitForTimeout(500);
    const done = await page.$('[data-testid="card-podman"] [data-testid="done-badge"]');
    if (done) break;
}

const vmSection = await page.$('[data-testid="card-podman"] [data-testid="vm-internal-section"]');
if (vmSection) {
    console.log('PODMAN_VM_INTERNAL_SECTION: present');
    const row0 = await page.$('[data-testid="card-podman"] [data-testid="vm-internal-row-0"]');
    const label0 = await page.$('[data-testid="card-podman"] [data-testid="vm-internal-label-0"]');
    const size0 = await page.$('[data-testid="card-podman"] [data-testid="vm-internal-size-0"]');
    console.log(`ELEM vm-internal-row-0: ${row0 ? 'present' : 'MISSING'}`);
    console.log(`ELEM vm-internal-label-0: ${label0 ? (await label0.textContent()).trim() : 'MISSING'}`);
    console.log(`ELEM vm-internal-size-0: ${size0 ? (await size0.textContent()).trim() : 'MISSING'}`);
} else {
    console.log('PODMAN_VM_INTERNAL_SECTION: absent');
    console.log('PODMAN_VM_INTERNAL_GRACEFUL: true');
}

console.log('DONE');