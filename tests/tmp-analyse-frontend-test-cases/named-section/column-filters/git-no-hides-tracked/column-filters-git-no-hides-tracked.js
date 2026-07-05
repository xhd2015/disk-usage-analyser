const BASE = process.env.SERVER_URL || 'http://localhost:8080';

async function countTrackedRows() {
    const gitCells = await page.$$('[data-testid^="node-modules-git-"]');
    let tracked = 0;
    for (const cell of gitCells) {
        const checked = await cell.evaluate(node => {
            if (node.tagName.toLowerCase() === 'input') return node.checked;
            const input = node.querySelector('input[type="checkbox"]');
            return input ? input.checked : false;
        });
        if (checked) tracked++;
    }
    return { total: gitCells.length, tracked };
}

async function selectGitFilter(value) {
    const el = await page.$('[data-testid="node-modules-filter-git"]');
    if (!el) return false;

    const tag = await el.evaluate(node => node.tagName.toLowerCase());
    if (tag === 'select') {
        await el.selectOption(value);
        return true;
    }

    const option = await el.$(`[data-value="${value}"], [value="${value}"]`);
    if (option) {
        await option.click();
        return true;
    }

    const items = await el.$$('.ant-segmented-item, .ant-select-item');
    for (const item of items) {
        const text = await item.evaluate(node => (node.textContent || '').trim().toLowerCase());
        if (text === value) {
            await item.click();
            return true;
        }
    }
    return false;
}

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

const scanBtn = await page.$('[data-testid="node-modules-scan-btn"]');
await scanBtn.click();

for (let i = 0; i < 120; i++) {
    await page.waitForTimeout(500);
    const done = await page.$('[data-testid="node-modules-section"] [data-testid="node-modules-done-badge"]');
    if (done) break;
}

const empty = await page.$('[data-testid="node-modules-empty-state"]');
if (empty) {
    console.log('NAMED_EMPTY_STATE: present');
    console.log('DONE');
    process.exit(0);
}

const gitFilter = await page.$('[data-testid="node-modules-filter-git"]');
console.log(`ELEM node-modules-filter-git: ${gitFilter ? 'present' : 'MISSING'}`);
if (!gitFilter) {
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

const before = await countTrackedRows();
console.log(`COUNT git-cells-before: ${before.total}`);
console.log(`COUNT git-tracked-before: ${before.tracked}`);

if (before.total === 0) {
    console.log('NAMED_EMPTY_STATE: present');
    console.log('DONE');
    process.exit(0);
}

const untrackedBefore = before.total - before.tracked;
if (before.tracked === 0 || untrackedBefore === 0) {
    console.log('SKIP no-git-mix');
    console.log('DONE');
    process.exit(0);
}

const selected = await selectGitFilter('no');
if (!selected) {
    console.log('FAIL git-filter-select: could not set no');
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

await page.waitForTimeout(300);

const after = await countTrackedRows();
console.log(`COUNT git-cells-after: ${after.total}`);
console.log(`COUNT git-tracked-after: ${after.tracked}`);

if (after.tracked === 0 && after.total <= untrackedBefore) {
    console.log('CHECK git-no-filter: pass');
} else {
    process.exitCode = 1;
}

console.log('DONE');