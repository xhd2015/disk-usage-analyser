const BASE = process.env.SERVER_URL || 'http://localhost:8080';

async function check(selector, label) {
    const el = await page.$(selector);
    if (!el) {
        console.log(`ELEM ${label}: MISSING`);
        return;
    }
    const text = (await el.textContent()) || '';
    const visible = await el.isVisible();
    console.log(`ELEM ${label}: "${text.trim()}" visible=${visible}`);
}

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

// Check each software card
const softwareCats = ['go', 'npm', 'bun', 'yarn', 'pnpm', 'pip', 'cargo',
                      'ruby', 'docker', 'podman', 'nginx', 'gradle', 'maven',
                      'android', 'brew', 'xcode', 'composer'];

for (const cat of softwareCats) {
    const card = await page.$(`[data-testid="card-${cat}"]`);
    if (card) {
        await check(`[data-testid="card-${cat}"]`, `card-${cat}`);
        await check(`[data-testid="card-${cat}"] [data-testid="card-label"]`, `card-${cat}-label`);
        await check(`[data-testid="card-${cat}"] [data-testid="card-size"]`, `card-${cat}-size`);
        await check(`[data-testid="card-${cat}"] [data-testid="reboot-safe-badge"]`, `card-${cat}-reboot-safe`);
    } else {
        console.log(`NOT_DETECTED ${cat}`);
    }
}

console.log('DONE');
