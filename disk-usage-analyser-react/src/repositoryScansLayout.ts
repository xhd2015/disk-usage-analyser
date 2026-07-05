export const ONE_MIB = 1048576;
export const TEN_MIB = 10485760;

export interface WorktreeRepoRow {
    repoPath: string;
    repoName: string;
    size: number;
    sizeHuman: string;
    fileCount: number;
}

export interface WorktreeHit {
    repoPath: string;
    repoName: string;
    path: string;
    head: string;
    isMain: boolean;
    size: number;
    sizeHuman: string;
    fileCount: number;
}

export interface BinaryHit {
    path: string;
    size: number;
    sizeHuman: string;
    kind: string;
    typeDesc: string;
    repoPath: string;
    repoName: string;
}

export function sortWorktreeRepos(repos: WorktreeRepoRow[]): WorktreeRepoRow[] {
    return [...repos].sort((a, b) => b.size - a.size);
}

export function sortLinkedWorktrees(hits: WorktreeHit[]): WorktreeHit[] {
    return [...hits].sort((a, b) => b.size - a.size);
}

export function filterWorktreeRepos(
    repos: WorktreeRepoRow[],
    linkedByRepo: Map<string, WorktreeHit[]>,
    showUnder10M: boolean,
): { repos: WorktreeRepoRow[]; linkedByRepo: Map<string, WorktreeHit[]> } {
    const filteredLinked = new Map<string, WorktreeHit[]>();

    for (const [repoPath, hits] of linkedByRepo.entries()) {
        const visibleHits = showUnder10M
            ? [...hits]
            : hits.filter(hit => hit.size >= TEN_MIB);
        if (visibleHits.length > 0) {
            filteredLinked.set(repoPath, visibleHits);
        }
    }

    const filteredRepos = showUnder10M
        ? [...repos]
        : repos.filter(repo => repo.size >= TEN_MIB);

    return { repos: filteredRepos, linkedByRepo: filteredLinked };
}

export function sortBinaryRepos(byRepo: Map<string, BinaryHit[]>): [string, BinaryHit[]][] {
    const entries = Array.from(byRepo.entries()).map(([repoPath, hits]) => {
        const sortedHits = [...hits].sort((a, b) => b.size - a.size);
        const total = sortedHits.reduce((sum, hit) => sum + hit.size, 0);
        return { repoPath, hits: sortedHits, total };
    });
    entries.sort((a, b) => b.total - a.total);
    return entries.map(({ repoPath, hits }) => [repoPath, hits]);
}

export function filterBinaryRepos(
    byRepo: Map<string, BinaryHit[]>,
    showUnder1M: boolean,
): Map<string, BinaryHit[]> {
    const filtered = new Map<string, BinaryHit[]>();

    for (const [repoPath, hits] of byRepo.entries()) {
        const visibleHits = showUnder1M
            ? [...hits]
            : hits.filter(hit => hit.size >= ONE_MIB);

        if (showUnder1M || visibleHits.reduce((sum, hit) => sum + hit.size, 0) >= ONE_MIB) {
            filtered.set(repoPath, visibleHits);
        }
    }

    return filtered;
}

export interface NamedHit {
    path: string;
    name: string;
    size: number;
    sizeHuman: string;
    repoPath: string;
    repoName: string;
    packageManager?: string;
    hasPackageJson?: boolean;
    gitTracked?: boolean;
    pnpmSharedSize?: number;
    pnpmSharedHuman?: string;
    bunSharedSize?: number;
    bunSharedHuman?: string;
    sharedSize?: number;
    sharedHuman?: string;
    enrichmentStatus?: 'pending' | 'resolved';
}

export function sortNamedRepos(byRepo: Map<string, NamedHit[]>): [string, NamedHit[]][] {
    const entries = Array.from(byRepo.entries()).map(([repoPath, hits]) => {
        const sortedHits = [...hits].sort((a, b) => b.size - a.size);
        const total = sortedHits.reduce((sum, hit) => sum + hit.size, 0);
        return { repoPath, hits: sortedHits, total };
    });
    entries.sort((a, b) => b.total - a.total);
    return entries.map(({ repoPath, hits }) => [repoPath, hits]);
}

export function sortNamedHits(hits: NamedHit[]): NamedHit[] {
    return [...hits].sort((a, b) => b.size - a.size);
}

export function filterNamedRepos(
    byRepo: Map<string, NamedHit[]>,
    showUnder1M: boolean,
): Map<string, NamedHit[]> {
    const filtered = new Map<string, NamedHit[]>();

    for (const [repoPath, hits] of byRepo.entries()) {
        const visibleHits = showUnder1M
            ? [...hits]
            : hits.filter(hit => hit.size >= ONE_MIB);

        if (showUnder1M || visibleHits.reduce((sum, hit) => sum + hit.size, 0) >= ONE_MIB) {
            filtered.set(repoPath, visibleHits);
        }
    }

    return filtered;
}

export type TriStateFilter = 'all' | 'yes' | 'no';
export type PmFilter = 'all' | 'npm' | 'pnpm' | 'yarn' | 'bun' | 'unknown';

export interface NamedColumnFilters {
    git: TriStateFilter;
    packageJson: TriStateFilter;
    pm: PmFilter;
}

export const defaultNamedColumnFilters: NamedColumnFilters = {
    git: 'all',
    packageJson: 'all',
    pm: 'all',
};

function matchesTriState(value: boolean | undefined, filter: TriStateFilter): boolean {
    if (filter === 'all') {
        return true;
    }
    if (filter === 'yes') {
        return value === true;
    }
    return value !== true;
}

function matchesPm(packageManager: string | undefined, filter: PmFilter): boolean {
    if (filter === 'all') {
        return true;
    }
    return (packageManager || 'unknown') === filter;
}

export function filterNamedReposByColumnFilters(
    byRepo: Map<string, NamedHit[]>,
    filters: NamedColumnFilters,
): Map<string, NamedHit[]> {
    const filtered = new Map<string, NamedHit[]>();

    for (const [repoPath, hits] of byRepo.entries()) {
        const visibleHits = hits.filter(hit =>
            matchesTriState(hit.gitTracked, filters.git)
            && matchesTriState(hit.hasPackageJson, filters.packageJson)
            && matchesPm(hit.packageManager, filters.pm),
        );

        if (visibleHits.length > 0) {
            filtered.set(repoPath, visibleHits);
        }
    }

    return filtered;
}