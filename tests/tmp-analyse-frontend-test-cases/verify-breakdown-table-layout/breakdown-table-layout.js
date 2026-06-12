const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

// Check Go card breakdown row at index 1 (extra path, now unified)
const goRow = await page.$('[data-testid="card-go"] [data-testid="breakdown-row-1"]');
if (!goRow) {
    console.log('ELEM go-breakdown-row-1: MISSING');
} else {
    console.log('ELEM go-breakdown-row-1: EXISTS');
    const goDisplay = await goRow.evaluate(el => getComputedStyle(el).display);
    const goJustify = await goRow.evaluate(el => getComputedStyle(el).justifyContent);
    console.log(`ROW_STYLE go-display: ${goDisplay}`);
    console.log(`ROW_STYLE go-justify: ${goJustify}`);

    const goLabel = await goRow.$('[data-testid="breakdown-label-1"]');
    const goSize = await goRow.$('[data-testid="breakdown-size-1"]');
    console.log(`ROW_CHILDREN go-has-label: ${goLabel !== null}`);
    console.log(`ROW_CHILDREN go-has-size: ${goSize !== null}`);
}

// Check Xcode card breakdown row at index 1
const xcRow = await page.$('[data-testid="card-xcode"] [data-testid="breakdown-row-1"]');
if (!xcRow) {
    console.log('ELEM xcode-breakdown-row-1: MISSING');
} else {
    console.log('ELEM xcode-breakdown-row-1: EXISTS');
    const xcDisplay = await xcRow.evaluate(el => getComputedStyle(el).display);
    const xcJustify = await xcRow.evaluate(el => getComputedStyle(el).justifyContent);
    console.log(`ROW_STYLE xcode-display: ${xcDisplay}`);
    console.log(`ROW_STYLE xcode-justify: ${xcJustify}`);

    const xcLabel = await xcRow.$('[data-testid="breakdown-label-1"]');
    const xcSize = await xcRow.$('[data-testid="breakdown-size-1"]');
    console.log(`ROW_CHILDREN xcode-has-label: ${xcLabel !== null}`);
    console.log(`ROW_CHILDREN xcode-has-size: ${xcSize !== null}`);
}

console.log('DONE');
