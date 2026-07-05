/**
 * Invoked by doctest Run() via: npx --yes tsx column-filters-harness.ts --op <op> --fixture <path>
 * Imports repositoryScansLayout.ts — RED until filterNamedReposByColumnFilters exists.
 */
import * as fs from 'node:fs';

import {
    filterNamedReposByColumnFilters,
    type NamedHit,
} from '../../../../../disk-usage-analyser-react/src/repositoryScansLayout.ts';

interface NamedColumnFilters {
    git: 'all' | 'yes' | 'no';
    packageJson: 'all' | 'yes' | 'no';
    pm: 'all' | 'npm' | 'pnpm' | 'yarn' | 'bun' | 'unknown';
}

interface Fixture {
    namedByRepo: Record<string, NamedHit[]>;
    filters: NamedColumnFilters;
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
        console.error('usage: column-filters-harness.ts --op <op> --fixture <path>');
        process.exit(2);
    }
    return { op, fixture };
}

function mapFromRecord(rec: Record<string, NamedHit[]>): Map<string, NamedHit[]> {
    const m = new Map<string, NamedHit[]>();
    for (const [k, v] of Object.entries(rec)) {
        m.set(k, v);
    }
    return m;
}

function runFilterNamedColumnFilters(fx: Fixture): Record<string, unknown> {
    const byRepo = mapFromRecord(fx.namedByRepo);
    const filtered = filterNamedReposByColumnFilters(byRepo, fx.filters);
    const repoOrder = Array.from(filtered.keys());
    const visiblePaths: string[] = [];
    for (const [, hits] of filtered.entries()) {
        for (const hit of hits) {
            visiblePaths.push(hit.path);
        }
    }
    return {
        repoOrder,
        visiblePaths,
        visibleCount: visiblePaths.length,
    };
}

function main() {
    const { op, fixture } = parseArgs(process.argv.slice(2));
    const raw = fs.readFileSync(fixture, 'utf8');
    const fx = JSON.parse(raw) as Fixture;

    if (op !== 'filter-named-column-filters') {
        console.error(`unknown op: ${op}`);
        process.exit(2);
    }

    const result = runFilterNamedColumnFilters(fx);
    process.stdout.write(JSON.stringify(result));
}

main();