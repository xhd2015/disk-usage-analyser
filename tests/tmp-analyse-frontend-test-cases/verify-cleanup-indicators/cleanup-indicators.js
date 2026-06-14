const BASE = process.env.SERVER_URL || 'http://localhost:8080';

async function check(selector, label) {
    const el = await page.$(selector);
    if (!el) {
        console.log(`ELEM ${label}: MISSING`);
        return;
    }
    const visible = await el.isVisible();
    console.log(`ELEM ${label}: visible=${visible}`);
}

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

const allCategories = [
    'trash', 'temp', 'cache', 'log',
    'go', 'npm', 'bun', 'yarn', 'pnpm', 'pip', 'cargo',
    'ruby', 'docker', 'podman', 'nginx', 'gradle', 'maven',
    'android', 'brew', 'xcode', 'composer',
];

for (const cat of allCategories) {
    await check(`[data-testid="card-${cat}"] [data-testid="cleanup-indicator"]`, `cleanup-indicator-${cat}`);
}

console.log('DONE');
