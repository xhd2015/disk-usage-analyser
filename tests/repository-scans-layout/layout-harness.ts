/**
 * Invoked by doctest Run() via: npx --yes tsx layout-harness.ts --op <op> --fixture <path>
 * Imports repositoryScansLayout.ts — RED until implementation exists.
 */
import * as fs from 'node:fs';

import {
    filterBinaryRepos,
    filterWorktreeRepos,
    sortBinaryRepos,
    sortLinkedWorktrees,
    sortWorktreeRepos,
} from '../../disk-usage-analyser-react/src/repositoryScansLayout.ts';

interface WorktreeRepoRow {
    repoPath: string;
    repoName: string;
    size: number;
    sizeHuman: string;
    fileCount: number;
}

interface WorktreeHit {
    repoPath: string;
    repoName: string;
    path: string;
    head: string;
    isMain: boolean;
    size: number;
    sizeHuman: string;
    fileCount: number;
}

interface BinaryHit {
    path: string;
    size: number;
    sizeHuman: string;
    kind: string;
    typeDesc: string;
    repoPath: string;
    repoName: string;
}

interface Fixture {
    worktreeRepos?: WorktreeRepoRow[];
    linkedByRepo?: Record<string, WorktreeHit[]>;
    binaryByRepo?: Record<string, BinaryHit[]>;
    showUnder1M?: boolean;
    showUnder10M?: boolean;
    insert?: WorktreeRepoRow;
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
        console.error('usage: layout-harness.ts --op <op> --fixture <path>');
        process.exit(2);
    }
    return { op, fixture };
}

function mapFromRecord<T>(rec: Record<string, T[]> | undefined): Map<string, T[]> {
    const m = new Map<string, T[]>();
    if (!rec) return m;
    for (const [k, v] of Object.entries(rec)) {
        m.set(k, v);
    }
    return m;
}

function linkedMapFromRecord(rec: Record<string, WorktreeHit[]> | undefined): Map<string, WorktreeHit[]> {
    return mapFromRecord(rec);
}

function binaryMapFromRecord(rec: Record<string, BinaryHit[]> | undefined): Map<string, BinaryHit[]> {
    return mapFromRecord(rec);
}

function runOp(op: string, fx: Fixture): Record<string, unknown> {
    switch (op) {
        case 'sort-worktree-repos': {
            const sorted = sortWorktreeRepos(fx.worktreeRepos ?? []);
            return {
                worktreeRepoOrder: sorted.map(r => r.repoPath),
                worktreeRepoSizes: sorted.map(r => r.size),
            };
        }
        case 'sort-binary-repos': {
            const byRepo = binaryMapFromRecord(fx.binaryByRepo);
            const sorted = sortBinaryRepos(byRepo);
            return {
                binaryRepoOrder: sorted.map(([repoPath]) => repoPath),
                binaryRepoTotals: sorted.map(([, hits]) => hits.reduce((s, h) => s + h.size, 0)),
            };
        }
        case 'sort-linked-worktrees': {
            const linked = fx.linkedByRepo ?? {};
            const out: Record<string, string[]> = {};
            const sizeOut: Record<string, number[]> = {};
            for (const [repo, hits] of Object.entries(linked)) {
                const sorted = sortLinkedWorktrees(hits);
                out[repo] = sorted.map(h => h.path);
                sizeOut[repo] = sorted.map(h => h.size);
            }
            return { linkedOrder: out, linkedSizes: sizeOut };
        }
        case 'filter-binary-repos': {
            const byRepo = binaryMapFromRecord(fx.binaryByRepo);
            const showUnder1M = fx.showUnder1M ?? false;
            const filtered = filterBinaryRepos(byRepo, showUnder1M);
            const rows: { repoPath: string; path: string; size: number }[] = [];
            for (const [repoPath, hits] of filtered.entries()) {
                for (const h of hits) {
                    rows.push({ repoPath, path: h.path, size: h.size });
                }
            }
            return {
                binaryRepoOrder: Array.from(filtered.keys()),
                visibleBinaryCount: rows.length,
                visibleBinaries: rows,
            };
        }
        case 'filter-worktree-repos': {
            const repos = fx.worktreeRepos ?? [];
            const linked = linkedMapFromRecord(fx.linkedByRepo);
            const showUnder10M = fx.showUnder10M ?? false;
            const { repos: filteredRepos, linkedByRepo: filteredLinked } = filterWorktreeRepos(
                repos,
                linked,
                showUnder10M,
            );
            const linkedOrder: Record<string, string[]> = {};
            for (const [repo, hits] of filteredLinked.entries()) {
                linkedOrder[repo] = hits.map(h => h.path);
            }
            return {
                worktreeRepoOrder: filteredRepos.map(r => r.repoPath),
                linkedOrder,
                visibleWorktreeRepoCount: filteredRepos.length,
            };
        }
        case 'resort-worktree-repos': {
            let repos = sortWorktreeRepos(fx.worktreeRepos ?? []);
            if (fx.insert) {
                repos = sortWorktreeRepos([...repos, fx.insert]);
            }
            return {
                worktreeRepoOrder: repos.map(r => r.repoPath),
                worktreeRepoSizes: repos.map(r => r.size),
            };
        }
        default:
            console.error(`unknown op: ${op}`);
            process.exit(2);
    }
}

const { op, fixture } = parseArgs(process.argv.slice(2));
const raw = fs.readFileSync(fixture, 'utf8');
const fx = JSON.parse(raw) as Fixture;
const result = runOp(op, fx);
process.stdout.write(JSON.stringify(result));