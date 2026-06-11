const BASE = process.env.SERVER_URL || 'http://localhost:8080';

// Go to home page
await page.goto(`${BASE}/`, { waitUntil: 'load' });
await page.waitForTimeout(500);
console.log(`URL home: ${page.url()}`);

// Find nav link with text "Tmp Files"
const links = await page.$$('nav a');
let tmpLink = null;
for (const link of links) {
    const text = await link.textContent();
    console.log(`NAV_LINK: "${text.trim()}"`);
    if (text && text.trim() === 'Tmp Files') {
        tmpLink = link;
    }
}

if (!tmpLink) {
    console.log('NAV_LINK Tmp Files: MISSING');
    process.exitCode = 1;
} else {
    console.log('NAV_LINK Tmp Files: FOUND');

    // Click the link
    await tmpLink.click();
    await page.waitForTimeout(500);
    console.log(`URL after click: ${page.url()}`);

    // Verify we're on /tmp-analyse
    if (page.url().includes('/tmp-analyse')) {
        console.log('URL tmp-analyse: REACHED');
    } else {
        console.log('URL tmp-analyse: NOT_REACHED');
        process.exitCode = 1;
    }

    // Verify the page heading exists
    const heading = await page.$('[data-testid="page-heading"]');
    if (heading) {
        const text = await heading.textContent();
        console.log(`ELEM page-heading: "${text.trim()}"`);
    } else {
        console.log('ELEM page-heading: MISSING');
        process.exitCode = 1;
    }
}

console.log('DONE');
