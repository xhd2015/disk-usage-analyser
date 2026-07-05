const BASE = process.env.SERVER_URL || 'http://localhost:8080';

async function readPmCells() {
    return page.$$eval('[data-testid^="node-modules-pkgmgr-"]', els =>
        els.map(el => (el.textContent || '').trim().toLowerCase()).filter(Boolean)
    ).catch(() => []);
}

async function selectPmFilter(value) {
    const el = await page.$('[data-testid="node-modules-filter-pm"]');
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

const pmFilter = await page.$('[data-testid="node-modules-filter-pm"]');
console.log(`ELEM node-modules-filter-pm: ${pmFilter ? 'present' : 'MISSING'}`);
if (!pmFilter) {
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

const pmsBefore = await readPmCells();
console.log(`COUNT pm-cells-before: ${pmsBefore.length}`);
const uniqueBefore = [...new Set(pmsBefore)];
console.log(`PM_SET before: ${uniqueBefore.join(',')}`);

if (pmsBefore.length === 0) {
    console.log('NAMED_EMPTY_STATE: present');
    console.log('DONE');
    process.exit(0);
}

const hasNpm = pmsBefore.includes('npm');
const hasOther = pmsBefore.some(pm => pm !== 'npm');
if (!hasNpm || !hasOther) {
    console.log('SKIP no-pm-mix');
    console.log('DONE');
    process.exit(0);
}

const selected = await selectPmFilter('npm');
if (!selected) {
    console.log('FAIL pm-filter-select: could not set npm');
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

await page.waitForTimeout(300);

const pmsAfter = await readPmCells();
console.log(`COUNT pm-cells-after: ${pmsAfter.length}`);
console.log(`PM_SET after: ${[...new Set(pmsAfter)].join(',')}`);

const allNpm = pmsAfter.length > 0 && pmsAfter.every(pm => pm === 'npm');
if (allNpm) {
    console.log('CHECK pm-filter-npm: pass');
} else {
    process.exitCode = 1;
}

console.log('DONE');