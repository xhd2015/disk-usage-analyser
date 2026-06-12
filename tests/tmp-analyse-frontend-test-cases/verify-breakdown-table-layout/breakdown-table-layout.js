const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

// Check Go card breakdown row
const goRow = await page.$('[data-testid="card-go"] [data-testid="extra-breakdown-row-0"]');
if (!goRow) {
    console.log('ELEM go-breakdown-row-0: MISSING');
} else {
    console.log('ELEM go-breakdown-row-0: EXISTS');
    // Check flexbox layout (display: flex, justify-content: space-between)
    const goDisplay = await goRow.evaluate(el => getComputedStyle(el).display);
    const goJustify = await goRow.evaluate(el => getComputedStyle(el).justifyContent);
    console.log(`ROW_STYLE go-display: ${goDisplay}`);
    console.log(`ROW_STYLE go-justify: ${goJustify}`);

    // Check it contains both label and size
    const goLabel = await goRow.$('[data-testid="extra-breakdown-label-0"]');
    const goSize = await goRow.$('[data-testid="extra-breakdown-size-0"]');
    console.log(`ROW_CHILDREN go-has-label: ${goLabel !== null}`);
    console.log(`ROW_CHILDREN go-has-size: ${goSize !== null}`);
}

// Check Xcode card breakdown row
const xcRow = await page.$('[data-testid="card-xcode"] [data-testid="extra-breakdown-row-0"]');
if (!xcRow) {
    console.log('ELEM xcode-breakdown-row-0: MISSING');
} else {
    console.log('ELEM xcode-breakdown-row-0: EXISTS');
    const xcDisplay = await xcRow.evaluate(el => getComputedStyle(el).display);
    const xcJustify = await xcRow.evaluate(el => getComputedStyle(el).justifyContent);
    console.log(`ROW_STYLE xcode-display: ${xcDisplay}`);
    console.log(`ROW_STYLE xcode-justify: ${xcJustify}`);

    const xcLabel = await xcRow.$('[data-testid="extra-breakdown-label-0"]');
    const xcSize = await xcRow.$('[data-testid="extra-breakdown-size-0"]');
    console.log(`ROW_CHILDREN xcode-has-label: ${xcLabel !== null}`);
    console.log(`ROW_CHILDREN xcode-has-size: ${xcSize !== null}`);
}

console.log('DONE');
