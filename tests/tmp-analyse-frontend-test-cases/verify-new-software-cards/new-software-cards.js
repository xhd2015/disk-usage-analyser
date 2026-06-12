const BASE = process.env.SERVER_URL;

async function check(selector, label) {
    const el = await page.$(selector);
    if (!el) {
        console.log(`ELEM ${label}: MISSING`);
        return null;
    }
    const text = (await el.textContent()) || '';
    const visible = await el.isVisible();
    console.log(`ELEM ${label}: "${text.trim()}" visible=${visible}`);
    return el;
}

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

const newCats = ['opencode', 'claude', 'codex', 'cursor'];
const multiPathCats = ['opencode', 'claude', 'codex', 'cursor'];

for (const cat of newCats) {
    const card = await page.$(`[data-testid="card-${cat}"]`);
    if (!card) {
        console.log(`CARD_NOT_FOUND ${cat}: card element missing (not detected on this system)`);
        continue;
    }
    console.log(`CARD_FOUND ${cat}: true`);

    // Card label
    await check(`[data-testid="card-${cat}"] [data-testid="card-label"]`, `card-${cat}-label`);

    // Card size
    await check(`[data-testid="card-${cat}"] [data-testid="card-size"]`, `card-${cat}-size`);

    // Reboot safe badge (all software caches are reboot-safe)
    const badge = await page.$(`[data-testid="card-${cat}"] [data-testid="reboot-safe-badge"]`);
    if (badge) {
        const badgeText = await badge.textContent();
        console.log(`REBOOT_SAFE ${cat}: "${badgeText}"`);
    } else {
        console.log(`REBOOT_SAFE ${cat}: MISSING`);
    }

    // Check for breakdown-items (all 4 new tools have multi-path)
    if (multiPathCats.includes(cat)) {
        const bd = await page.$(`[data-testid="card-${cat}"] [data-testid="breakdown-items"]`);
        console.log(`HAS_BREAKDOWN ${cat}: ${bd !== null}`);
    }

    // Check for card-path (should NOT exist for multi-path tools)
    const cp = await page.$(`[data-testid="card-${cat}"] [data-testid="card-path"]`);
    console.log(`HAS_CARD_PATH ${cat}: ${cp !== null}`);
}

// OpenCode: verify specific breakdown rows exist
const opencodeCard = await page.$('[data-testid="card-opencode"]');
if (opencodeCard) {
    // Should have multiple breakdown rows
    for (let i = 0; i <= 3; i++) {
        const row = await page.$(`[data-testid="card-opencode"] [data-testid="breakdown-row-${i}"]`);
        console.log(`BREAKDOWN_ROW opencode-row-${i}: ${row !== null ? 'EXISTS' : 'MISSING'}`);
    }
}

// Claude: verify specific breakdown rows
const claudeCard = await page.$('[data-testid="card-claude"]');
if (claudeCard) {
    for (let i = 0; i <= 2; i++) {
        const row = await page.$(`[data-testid="card-claude"] [data-testid="breakdown-row-${i}"]`);
        console.log(`BREAKDOWN_ROW claude-row-${i}: ${row !== null ? 'EXISTS' : 'MISSING'}`);
    }
}

console.log('DONE');
