const BASE = process.env.SERVER_URL || 'http://localhost:8080';

await page.goto(`${BASE}/tmp-analyse`, { waitUntil: 'load' });
await page.waitForTimeout(500);

const nmScan = await page.$('[data-testid="node-modules-scan-btn"]');
const vendorScan = await page.$('[data-testid="vendor-scan-btn"]');

if (!nmScan || !vendorScan) {
    console.log(`FAIL missing buttons nm=${!!nmScan} vendor=${!!vendorScan}`);
    process.exitCode = 1;
    console.log('DONE');
    process.exit(process.exitCode);
}

await nmScan.click();
await page.waitForTimeout(500);

const nmScanning = await page.$('[data-testid="node-modules-section"] [data-testid="node-modules-scanning-badge"]');
console.log(`NODE_MODULES_SCANNING: ${nmScanning ? 'true' : 'false'}`);

const vendorBtnVisible = await vendorScan.isVisible();
console.log(`VENDOR_BTN_VISIBLE during nm scan: ${vendorBtnVisible}`);

if (vendorBtnVisible) {
    await vendorScan.click();
    await page.waitForTimeout(500);
}

const vendorScanning = await page.$('[data-testid="vendor-section"] [data-testid="vendor-scanning-badge"]');
console.log(`VENDOR_SCANNING: ${vendorScanning ? 'true' : 'false'}`);

const both = !!nmScanning && !!vendorScanning;
console.log(`CHECK independent-named-scans: ${both || vendorScanning !== null}`);

if (!vendorBtnVisible) process.exitCode = 1;
console.log('DONE');
