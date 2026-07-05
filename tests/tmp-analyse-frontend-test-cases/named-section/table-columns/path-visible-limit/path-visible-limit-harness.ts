/**
 * Invoked by doctest Run() via: npx --yes tsx path-visible-limit-harness.ts --op path-visible-limit
 * Imports PATH_VISIBLE_CHAR_LIMIT from pathDisplay.ts — RED until limit is 56.
 */
import { PATH_VISIBLE_CHAR_LIMIT } from '../../../../../disk-usage-analyser-react/src/pathDisplay.ts';

function parseArgs(argv: string[]): { op: string } {
    let op = '';
    for (let i = 0; i < argv.length; i++) {
        if (argv[i] === '--op' && argv[i + 1]) {
            op = argv[++i];
        }
    }
    if (!op) {
        console.error('usage: path-visible-limit-harness.ts --op path-visible-limit');
        process.exit(2);
    }
    return { op };
}

function main() {
    const { op } = parseArgs(process.argv.slice(2));
    if (op !== 'path-visible-limit') {
        console.error(`unknown op: ${op}`);
        process.exit(2);
    }

    process.stdout.write(JSON.stringify({
        limit: PATH_VISIBLE_CHAR_LIMIT,
    }));
}

main();