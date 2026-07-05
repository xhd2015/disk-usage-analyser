/**
 * Invoked by doctest Run() via: npx --yes tsx path-display-harness.ts --op <op> --fixture <path>
 * Imports pathDisplay.ts — RED until implementation exists.
 */
import * as fs from 'node:fs';

import { truncatePathKeepSuffix } from '../../../../../disk-usage-analyser-react/src/pathDisplay.ts';

interface Fixture {
    path: string;
    maxVisibleChars: number;
}

function parseArgs(argv: string[]): { op: string; fixture: string } {
    let op = '';
    let fixture = '';
    for (let i = 0; i < argv.length; i++) {
        if (argv[i] === '--op' && argv[i + 1]) {
            op = argv[++i];
        } else if (argv[i] === '--fixture' && argv[i + 1]) {
            fixture = argv[++i];
        }
    }
    if (!op || !fixture) {
        console.error('usage: path-display-harness.ts --op <op> --fixture <path>');
        process.exit(2);
    }
    return { op, fixture };
}

function main() {
    const { op, fixture } = parseArgs(process.argv.slice(2));
    const data = JSON.parse(fs.readFileSync(fixture, 'utf8')) as Fixture;

    if (op !== 'truncate-path') {
        console.error(`unknown op: ${op}`);
        process.exit(2);
    }

    const display = truncatePathKeepSuffix(data.path, data.maxVisibleChars);
    process.stdout.write(JSON.stringify({
        path: data.path,
        maxVisibleChars: data.maxVisibleChars,
        display,
    }));
}

main();