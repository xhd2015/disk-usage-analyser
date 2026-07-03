import { useState, useRef, useEffect } from 'react';
import { Button, Card, Space, Typography, Row, Col, Tag, Collapse, Alert, Popover, Modal, Checkbox } from 'antd';
import { PlayCircleOutlined, StopOutlined, SyncOutlined, CheckCircleOutlined, ClockCircleOutlined, QuestionCircleOutlined } from '@ant-design/icons';
import {
    filterBinaryRepos,
    filterWorktreeRepos,
    sortBinaryRepos,
    sortLinkedWorktrees,
    sortWorktreeRepos,
    filterNamedRepos,
    sortNamedRepos,
    sortNamedHits,
} from './repositoryScansLayout';
import type { NamedHit } from './repositoryScansLayout';

const { Text } = Typography;

interface TmpRuntimeItem {
    type: string;
    totalCount: number;
    activeCount: number;
    size: number;
    reclaimable?: number;
}

interface TmpVmStorageItem {
    label: string;
    path: string;
    size: number;
}

interface TmpVmInternal {
    items: TmpVmStorageItem[];
    totalSize: number;
    machineRunning: boolean;
}

interface TmpLocation {
    label: string;
    category: string;
    size: number;
    fileCount: number;
    rebootSafe: boolean;
    reclaimable: boolean;
    detected: boolean;
    breakdownItems?: { path: string; size: number; fileCount: number }[];
    runtimeItems?: TmpRuntimeItem[];
    vmInternal?: TmpVmInternal;
}

interface CleanupSuggestion {
    command: string;
    removes: string;
    recoverable: string;
}

function formatBytes(bytes: number): string {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}

const cleanupSuggestions: Record<string, CleanupSuggestion[]> = {
    trash: [
        { command: 'rm -rf ~/.Trash/*', removes: 'Everything in Trash', recoverable: 'Permanently deleted; cannot recover' },
    ],
    cache: [
        { command: 'rm -rf ~/Library/Caches/*', removes: 'All application caches', recoverable: 'Yes, apps recreate on next launch' },
    ],
    log: [
        { command: 'rm -rf ~/Library/Logs/*', removes: 'All user log files', recoverable: 'Yes, apps recreate log files automatically' },
    ],
    temp: [
        { command: 'rm -rf /tmp/*', removes: 'All temporary files', recoverable: 'Cleared on reboot; safe to delete' },
    ],
    go: [
        { command: 'go clean -cache', removes: '~/Library/Caches/go-build', recoverable: 'Yes, rebuilt automatically on next go build' },
        { command: 'go clean -modcache', removes: '~/go/pkg/mod', recoverable: 'Yes, re-downloaded automatically on next build' },
    ],
    npm: [
        { command: 'npm cache clean --force', removes: '~/.npm/_cacache/', recoverable: 'Yes, re-downloaded automatically on npm install' },
        { command: 'npm cache verify', removes: 'Nothing (verifies cache integrity only)', recoverable: 'No files removed' },
    ],
    bun: [
        { command: 'rm -rf ~/.bun/install/cache', removes: 'Bun package cache', recoverable: 'Yes, re-downloaded on bun install' },
    ],
    yarn: [
        { command: 'yarn cache clean', removes: 'Yarn package cache', recoverable: 'Yes, re-downloaded on yarn install' },
    ],
    pnpm: [
        { command: 'pnpm store prune', removes: 'Unused packages from pnpm store', recoverable: 'Yes, re-downloaded if needed' },
    ],
    pip: [
        { command: 'pip cache purge', removes: 'pip download cache', recoverable: 'Yes, re-downloaded on pip install' },
    ],
    cargo: [
        { command: 'rm -rf ~/.cargo/registry/cache', removes: 'Rust crate cache', recoverable: 'Yes, re-downloaded on cargo build' },
    ],
    ruby: [
        { command: 'gem cleanup', removes: 'Old gem versions', recoverable: 'Yes, current versions kept; old can be reinstalled' },
    ],
    docker: [
        { command: 'docker system prune -a', removes: 'Unused images, build cache, volumes', recoverable: 'Yes, re-pulled/rebuilt when needed' },
    ],
    podman: [
        { command: 'podman machine ssh', removes: 'Opens the Podman VM shell — run all commands below as user core, not root (root sees a separate empty store)', recoverable: 'No files removed' },
        { command: 'podman system df -v', removes: 'Nothing (shows images, containers, volumes, and build cache usage)', recoverable: 'No files removed' },
        { command: 'podman ps -a --size', removes: 'Nothing (lists containers with writable layer sizes)', recoverable: 'No files removed' },
        { command: 'podman images', removes: 'Nothing (lists images, including dangling <none> build leftovers)', recoverable: 'No files removed' },
        { command: 'du -sh ~/.local/share/containers/storage/overlay', removes: 'Nothing (raw overlay directory size on disk)', recoverable: 'No files removed' },
        { command: 'podman image prune -f', removes: 'Dangling <none> images only', recoverable: 'Yes, re-pulled when needed' },
        { command: 'podman container prune -f', removes: 'All stopped containers (running containers are kept)', recoverable: 'Yes, containers must be recreated if needed' },
        { command: 'podman image prune -a -f', removes: 'All unused images (dangling + unused tagged builds)', recoverable: 'Yes, re-pulled when needed; images referenced by running containers are kept' },
        { command: 'podman image prune -a --filter until=672h -f', removes: 'Unused images older than 4 weeks', recoverable: 'Yes, re-pulled when needed' },
        { command: 'podman rmi <image_id>', removes: 'A specific image by ID or name', recoverable: 'Yes, re-pulled when needed; run podman rm <container_id> first if image is in use' },
        { command: 'podman system prune -a -f --volumes', removes: 'Stopped containers, unused images, networks, and volumes', recoverable: 'Yes, re-pulled/rebuilt when needed' },
        { command: 'podman system reset -f', removes: 'All Podman data for user core (containers, images, volumes, orphaned overlay dirs)', recoverable: 'No — must re-pull images and recreate containers; use when prune leaves huge overlay/ or metadata is corrupted' },
    ],
    nginx: [
        { command: 'rm -f /usr/local/var/log/nginx/*.log', removes: 'Nginx access and error logs', recoverable: 'Yes, log files auto-created on next access' },
    ],
    gradle: [
        { command: 'rm -rf ~/.gradle/caches', removes: 'Gradle build cache', recoverable: 'Yes, rebuilt on next Gradle build' },
    ],
    maven: [
        { command: 'rm -rf ~/.m2/repository', removes: 'All downloaded Maven artifacts', recoverable: 'Yes, re-downloaded on next Maven build' },
    ],
    android: [
        { command: 'rm -rf ~/Library/Android/sdk/.temp', removes: 'Android SDK temporary files', recoverable: 'Yes, recreated if needed' },
    ],
    brew: [
        { command: 'brew cleanup', removes: 'Old formula versions', recoverable: 'Yes, current versions kept; re-downloadable' },
    ],
    xcode: [
        { command: 'rm -rf ~/Library/Developer/Xcode/DerivedData', removes: 'Xcode build products', recoverable: 'Yes, rebuilt on next Xcode build' },
        { command: 'xcrun simctl shutdown all && xcrun simctl delete all', removes: 'All iOS Simulators (devices)', recoverable: 'Yes, recreated via Xcode > Settings > Devices' },
    ],
    composer: [
        { command: 'composer clear-cache', removes: 'PHP Composer package cache', recoverable: 'Yes, re-downloaded on composer install' },
    ],
    opencode: [
        { command: 'rm -rf ~/.local/share/opencode/snapshot ~/.local/share/opencode/project ~/.local/share/opencode/tool-output ~/.local/share/opencode/storage ~/.local/share/opencode/log ~/.cache/opencode ~/.local/state/opencode', removes: 'OpenCode caches, snapshots, logs', recoverable: 'Some auto-recreated; snapshots lost permanently' },
    ],
    claude: [
        { command: 'rm -rf ~/.claude/cache', removes: 'Claude Code cache', recoverable: 'Yes, recreated on next use' },
    ],
    codex: [
        { command: 'rm -rf ~/.codex ~/Library/Application\\ Support/codex', removes: 'Codex (OpenAI) caches', recoverable: 'Some auto-recreated; settings may be lost' },
    ],
    cursor: [
        { command: 'rm -rf ~/Library/Application\\ Support/Cursor ~/Library/Application\\ Support/Caches/cursor-updater ~/Library/Caches/cursor-compile-cache', removes: 'Cursor caches and updater files', recoverable: 'Some auto-recreated; settings may be lost' },
    ],
};

const categoryColors: Record<string, string> = {
    trash: 'red',
    temp: 'orange',
    cache: 'blue',
    log: 'purple',
    swap: 'default',
    go: 'cyan',
    npm: 'green',
    bun: 'lime',
    yarn: 'geekblue',
    pnpm: 'gold',
    pip: 'magenta',
    cargo: 'volcano',
    ruby: 'red',
    docker: 'blue',
    podman: 'purple',
    nginx: 'green',
    gradle: 'orange',
    maven: 'cyan',
    android: 'lime',
    brew: 'gold',
    xcode: 'geekblue',
    composer: 'magenta',
};

type LocStatus = 'idle' | 'pending' | 'scanning' | 'done';

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

function testIdKey(value: string): string {
    return value.replace(/[^a-zA-Z0-9._-]/g, '_');
}

function TmpFilesAnalyse() {
    const [scanning, setScanning] = useState(false);
    const [locations, setLocations] = useState<TmpLocation[]>([]);
    const [locStatuses, setLocStatuses] = useState<Record<string, LocStatus>>({});
    const eventSourceRef = useRef<EventSource | null>(null);
    const mainScanActiveRef = useRef(false);
    const [totalSize, setTotalSize] = useState(0);
    const [reclaimableSize, setReclaimableSize] = useState(0);
    const [platformError, setPlatformError] = useState<string | null>(null);

    const [worktreesScanning, setWorktreesScanning] = useState(false);
    const [worktreesDone, setWorktreesDone] = useState(false);
    const [worktreeRepos, setWorktreeRepos] = useState<WorktreeRepoRow[]>([]);
    const [worktreeHits, setWorktreeHits] = useState<WorktreeHit[]>([]);
    const worktreesESRef = useRef<EventSource | null>(null);

    const [binariesScanning, setBinariesScanning] = useState(false);
    const [binariesDone, setBinariesDone] = useState(false);
    const [binaryHits, setBinaryHits] = useState<BinaryHit[]>([]);
    const [showBinaryUnder1M, setShowBinaryUnder1M] = useState(false);
    const [showWorktreeUnder10M, setShowWorktreeUnder10M] = useState(false);
    const [selectedBinaryPaths, setSelectedBinaryPaths] = useState<Set<string>>(new Set());
    const [deleteModalOpen, setDeleteModalOpen] = useState(false);
    const [deleteSuccess, setDeleteSuccess] = useState(false);
    const [deleting, setDeleting] = useState(false);
    const binariesESRef = useRef<EventSource | null>(null);

    const [namedScanning, setNamedScanning] = useState(false);
    const [namedDone, setNamedDone] = useState(false);
    const [namedHits, setNamedHits] = useState<NamedHit[]>([]);
    const [showNamedUnder1M, setShowNamedUnder1M] = useState(false);
    const [selectedNamedPaths, setSelectedNamedPaths] = useState<Set<string>>(new Set());
    const [namedDeleteModalOpen, setNamedDeleteModalOpen] = useState(false);
    const [namedDeleteSuccess, setNamedDeleteSuccess] = useState(false);
    const [namedDeleting, setNamedDeleting] = useState(false);
    const namedESRef = useRef<EventSource | null>(null);

    const [vendorScanning, setVendorScanning] = useState(false);
    const [vendorDone, setVendorDone] = useState(false);
    const [vendorHits, setVendorHits] = useState<NamedHit[]>([]);
    const [showVendorUnder1M, setShowVendorUnder1M] = useState(false);
    const [selectedVendorPaths, setSelectedVendorPaths] = useState<Set<string>>(new Set());
    const [vendorDeleteModalOpen, setVendorDeleteModalOpen] = useState(false);
    const [vendorDeleteSuccess, setVendorDeleteSuccess] = useState(false);
    const [vendorDeleting, setVendorDeleting] = useState(false);
    const vendorESRef = useRef<EventSource | null>(null);

    const [worktreeRunningTotal, setWorktreeRunningTotal] = useState(0);
    const [binaryRunningTotal, setBinaryRunningTotal] = useState(0);
    const [namedRunningTotal, setNamedRunningTotal] = useState(0);
    const [vendorRunningTotal, setVendorRunningTotal] = useState(0);

    useEffect(() => {
        fetch('/api/tmp-analyse-locations')
            .then(res => res.json())
            .then((locs: TmpLocation[]) => {
                setLocations(locs);
                const stat: Record<string, LocStatus> = {};
                locs.forEach(l => { stat[l.label] = 'idle'; });
                setLocStatuses(stat);
            })
            .catch(err => console.error('Failed to fetch locations:', err));
    }, []);

    const finalizeMainScanStatuses = () => {
        setLocStatuses(prev => {
            const next = { ...prev };
            for (const label of Object.keys(next)) {
                if (next[label] === 'pending' || next[label] === 'scanning') {
                    next[label] = 'done';
                }
            }
            return next;
        });
    };

    const startScan = () => {
        setScanning(true);
        setPlatformError(null);
        mainScanActiveRef.current = true;

        const pending: Record<string, LocStatus> = {};
        locations.forEach(l => { pending[l.label] = 'pending'; });
        setLocStatuses(pending);

        const es = new EventSource('/api/tmp-analyse');
        eventSourceRef.current = es;

        es.addEventListener('locations', (e) => {
            const locs: TmpLocation[] = JSON.parse((e as MessageEvent).data);
            setLocations(locs);
            const locStat: Record<string, LocStatus> = {};
            locs.forEach(l => { locStat[l.label] = 'idle'; });
            setLocStatuses(locStat);

            const pending2: Record<string, LocStatus> = {};
            locs.filter(l => l.detected).forEach(l => { pending2[l.label] = 'pending'; });
            setLocStatuses(prev => ({ ...prev, ...pending2 }));
        });

        es.addEventListener('unsupported_platform', (e) => {
            const d = JSON.parse((e as MessageEvent).data);
            setPlatformError(`This feature is only supported on macOS. Current OS: ${d.os}`);
            mainScanActiveRef.current = false;
            es.close();
            eventSourceRef.current = null;
            setScanning(false);
        });

        es.addEventListener('progress', (e) => {
            if (!mainScanActiveRef.current) return;
            const p = JSON.parse((e as MessageEvent).data);
            setLocations(prev => prev.map(loc => {
                if (loc.label !== p.label) return loc;
                const updated = { ...loc, size: p.size, fileCount: p.fileCount };
                if (p.breakdownIndex !== undefined && loc.breakdownItems) {
                    const items = [...loc.breakdownItems];
                    if (p.breakdownPath) {
                        const pathIdx = items.findIndex(item => item.path === p.breakdownPath);
                        if (pathIdx >= 0) {
                            items[pathIdx] = {
                                ...items[pathIdx],
                                size: p.breakdownSize,
                                fileCount: p.breakdownFileCount,
                            };
                        } else {
                            items.push({
                                path: p.breakdownPath,
                                size: p.breakdownSize,
                                fileCount: p.breakdownFileCount,
                            });
                        }
                    } else {
                        const idx = p.breakdownIndex as number;
                        if (idx < items.length) {
                            items[idx] = {
                                ...items[idx],
                                size: p.breakdownSize,
                                fileCount: p.breakdownFileCount,
                            };
                        }
                    }
                    updated.breakdownItems = items;
                }
                return updated;
            }));
            setLocStatuses(prev => {
                if (!mainScanActiveRef.current) return prev;
                return { ...prev, [p.label]: 'scanning' };
            });
            setTotalSize(p.totalSize);
            setReclaimableSize(p.reclaimableSize);
        });

        es.addEventListener('location', (e) => {
            if (!mainScanActiveRef.current) return;
            const loc: TmpLocation = JSON.parse((e as MessageEvent).data);
            setLocations(prev => prev.map(p =>
                p.label === loc.label
                    ? { ...p, size: loc.size, fileCount: loc.fileCount, breakdownItems: loc.breakdownItems, runtimeItems: loc.runtimeItems, vmInternal: loc.vmInternal }
                    : p
            ));
            setLocStatuses(prev => {
                if (!mainScanActiveRef.current) return prev;
                return { ...prev, [loc.label]: 'done' };
            });
        });

        es.addEventListener('done', () => {
            mainScanActiveRef.current = false;
            es.close();
            eventSourceRef.current = null;
            setScanning(false);
            finalizeMainScanStatuses();
        });

        es.addEventListener('server_error', (e) => {
            const d = JSON.parse((e as MessageEvent).data);
            console.error('SSE error:', d.error);
            mainScanActiveRef.current = false;
            es.close();
            eventSourceRef.current = null;
            setScanning(false);
            finalizeMainScanStatuses();
        });

        es.onerror = () => {
            if (es.readyState === EventSource.CLOSED) return;
            mainScanActiveRef.current = false;
            es.close();
            eventSourceRef.current = null;
            setScanning(false);
            finalizeMainScanStatuses();
        };
    };

    const stopScan = () => {
        mainScanActiveRef.current = false;
        if (eventSourceRef.current) {
            eventSourceRef.current.close();
            eventSourceRef.current = null;
        }
        setScanning(false);
    };

    const startWorktreesScan = () => {
        if (worktreesESRef.current) {
            worktreesESRef.current.close();
            worktreesESRef.current = null;
        }
        setWorktreesScanning(true);
        setWorktreesDone(false);
        setWorktreeRepos([]);
        setWorktreeHits([]);
        setWorktreeRunningTotal(0);

        const es = new EventSource('/api/tmp-worktrees-scan');
        worktreesESRef.current = es;

        es.addEventListener('repo', (e) => {
            const row: WorktreeRepoRow = JSON.parse((e as MessageEvent).data);
            setWorktreeRepos(prev => [...prev, row]);
            setWorktreeRunningTotal(prev => prev + row.size);
        });

        es.addEventListener('worktree', (e) => {
            const hit: WorktreeHit = JSON.parse((e as MessageEvent).data);
            setWorktreeHits(prev => [...prev, hit]);
            setWorktreeRunningTotal(prev => prev + hit.size);
        });

        es.addEventListener('done', () => {
            es.close();
            worktreesESRef.current = null;
            setWorktreesScanning(false);
            setWorktreesDone(true);
        });

        es.addEventListener('server_error', (e) => {
            console.error('Worktrees SSE error:', JSON.parse((e as MessageEvent).data));
            es.close();
            worktreesESRef.current = null;
            setWorktreesScanning(false);
        });

        es.onerror = () => {
            if (es.readyState === EventSource.CLOSED) return;
            es.close();
            worktreesESRef.current = null;
            setWorktreesScanning(false);
        };
    };

    const stopWorktreesScan = () => {
        if (worktreesESRef.current) {
            worktreesESRef.current.close();
            worktreesESRef.current = null;
        }
        setWorktreesScanning(false);
    };

    const startBinariesScan = () => {
        if (binariesESRef.current) {
            binariesESRef.current.close();
            binariesESRef.current = null;
        }
        setBinariesScanning(true);
        setBinariesDone(false);
        setBinaryHits([]);
        setSelectedBinaryPaths(new Set());
        setDeleteSuccess(false);
        setBinaryRunningTotal(0);

        const es = new EventSource('/api/tmp-binaries-scan');
        binariesESRef.current = es;

        es.addEventListener('binary', (e) => {
            const hit: BinaryHit = JSON.parse((e as MessageEvent).data);
            setBinaryHits(prev => [...prev, hit]);
            setBinaryRunningTotal(prev => prev + hit.size);
        });

        es.addEventListener('done', () => {
            es.close();
            binariesESRef.current = null;
            setBinariesScanning(false);
            setBinariesDone(true);
        });

        es.addEventListener('server_error', (e) => {
            console.error('Binaries SSE error:', JSON.parse((e as MessageEvent).data));
            es.close();
            binariesESRef.current = null;
            setBinariesScanning(false);
        });

        es.onerror = () => {
            if (es.readyState === EventSource.CLOSED) return;
            es.close();
            binariesESRef.current = null;
            setBinariesScanning(false);
        };
    };

    const stopBinariesScan = () => {
        if (binariesESRef.current) {
            binariesESRef.current.close();
            binariesESRef.current = null;
        }
        setBinariesScanning(false);
    };

    const startNamedScan = () => {
        if (namedESRef.current) {
            namedESRef.current.close();
            namedESRef.current = null;
        }
        setNamedScanning(true);
        setNamedDone(false);
        setNamedHits([]);
        setSelectedNamedPaths(new Set());
        setNamedDeleteSuccess(false);
        setNamedRunningTotal(0);

        const es = new EventSource('/api/tmp-named-scan?name=node_modules');
        namedESRef.current = es;

        es.addEventListener('named', (e) => {
            const hit: NamedHit = JSON.parse((e as MessageEvent).data);
            setNamedHits(prev => [...prev, hit]);
            setNamedRunningTotal(prev => prev + hit.size);
        });

        es.addEventListener('done', () => {
            es.close();
            namedESRef.current = null;
            setNamedScanning(false);
            setNamedDone(true);
        });

        es.addEventListener('server_error', (e) => {
            console.error('Named SSE error:', JSON.parse((e as MessageEvent).data));
            es.close();
            namedESRef.current = null;
            setNamedScanning(false);
        });

        es.onerror = () => {
            if (es.readyState === EventSource.CLOSED) return;
            es.close();
            namedESRef.current = null;
            setNamedScanning(false);
        };
    };

    const stopNamedScan = () => {
        if (namedESRef.current) {
            namedESRef.current.close();
            namedESRef.current = null;
        }
        setNamedScanning(false);
    };

    const toggleNamedPath = (path: string, checked: boolean) => {
        setSelectedNamedPaths(prev => {
            const next = new Set(prev);
            if (checked) {
                next.add(path);
            } else {
                next.delete(path);
            }
            return next;
        });
    };

    const toggleNamedRepo = (repoPath: string, checked: boolean) => {
        const paths = namedHits.filter(hit => hit.repoPath === repoPath).map(hit => hit.path);
        setSelectedNamedPaths(prev => {
            const next = new Set(prev);
            for (const path of paths) {
                if (checked) {
                    next.add(path);
                } else {
                    next.delete(path);
                }
            }
            return next;
        });
    };

    const confirmDeleteNamed = async () => {
        setNamedDeleting(true);
        try {
            const res = await fetch('/api/tmp-named-delete', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ paths: Array.from(selectedNamedPaths) }),
            });
            const result = await res.json();
            const deletedSet = new Set<string>(result.deleted || []);
            setNamedHits(prev => prev.filter(hit => !deletedSet.has(hit.path)));
            setSelectedNamedPaths(new Set());
            setNamedDeleteSuccess(true);
            setNamedDeleteModalOpen(false);
        } catch (err) {
            console.error('Named delete failed:', err);
        } finally {
            setNamedDeleting(false);
        }
    };

    const startVendorScan = () => {
        if (vendorESRef.current) {
            vendorESRef.current.close();
            vendorESRef.current = null;
        }
        setVendorScanning(true);
        setVendorDone(false);
        setVendorHits([]);
        setSelectedVendorPaths(new Set());
        setVendorDeleteSuccess(false);
        setVendorRunningTotal(0);

        const es = new EventSource('/api/tmp-named-scan?name=vendor');
        vendorESRef.current = es;

        es.addEventListener('named', (e) => {
            const hit: NamedHit = JSON.parse((e as MessageEvent).data);
            setVendorHits(prev => [...prev, hit]);
            setVendorRunningTotal(prev => prev + hit.size);
        });

        es.addEventListener('done', () => {
            es.close();
            vendorESRef.current = null;
            setVendorScanning(false);
            setVendorDone(true);
        });

        es.addEventListener('server_error', (e) => {
            console.error('Vendor SSE error:', JSON.parse((e as MessageEvent).data));
            es.close();
            vendorESRef.current = null;
            setVendorScanning(false);
        });

        es.onerror = () => {
            if (es.readyState === EventSource.CLOSED) return;
            es.close();
            vendorESRef.current = null;
            setVendorScanning(false);
        };
    };

    const stopVendorScan = () => {
        if (vendorESRef.current) {
            vendorESRef.current.close();
            vendorESRef.current = null;
        }
        setVendorScanning(false);
    };

    const toggleVendorPath = (path: string, checked: boolean) => {
        setSelectedVendorPaths(prev => {
            const next = new Set(prev);
            if (checked) {
                next.add(path);
            } else {
                next.delete(path);
            }
            return next;
        });
    };

    const toggleVendorRepo = (repoPath: string, checked: boolean) => {
        const paths = vendorHits.filter(hit => hit.repoPath === repoPath).map(hit => hit.path);
        setSelectedVendorPaths(prev => {
            const next = new Set(prev);
            for (const path of paths) {
                if (checked) {
                    next.add(path);
                } else {
                    next.delete(path);
                }
            }
            return next;
        });
    };

    const confirmDeleteVendor = async () => {
        setVendorDeleting(true);
        try {
            const res = await fetch('/api/tmp-named-delete', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ paths: Array.from(selectedVendorPaths) }),
            });
            const result = await res.json();
            const deletedSet = new Set<string>(result.deleted || []);
            setVendorHits(prev => prev.filter(hit => !deletedSet.has(hit.path)));
            setSelectedVendorPaths(new Set());
            setVendorDeleteSuccess(true);
            setVendorDeleteModalOpen(false);
        } catch (err) {
            console.error('Vendor delete failed:', err);
        } finally {
            setVendorDeleting(false);
        }
    };

    const selectedBinaryTotal = binaryHits
        .filter(hit => selectedBinaryPaths.has(hit.path))
        .reduce((sum, hit) => sum + hit.size, 0);

    const toggleBinaryPath = (path: string, checked: boolean) => {
        setSelectedBinaryPaths(prev => {
            const next = new Set(prev);
            if (checked) {
                next.add(path);
            } else {
                next.delete(path);
            }
            return next;
        });
    };

    const toggleRepoBinaries = (repoPath: string, checked: boolean) => {
        const paths = binaryHits.filter(hit => hit.repoPath === repoPath).map(hit => hit.path);
        setSelectedBinaryPaths(prev => {
            const next = new Set(prev);
            for (const path of paths) {
                if (checked) {
                    next.add(path);
                } else {
                    next.delete(path);
                }
            }
            return next;
        });
    };

    const confirmDeleteBinaries = async () => {
        setDeleting(true);
        try {
            const res = await fetch('/api/tmp-binaries-delete', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ paths: Array.from(selectedBinaryPaths) }),
            });
            const result = await res.json();
            const deletedSet = new Set<string>(result.deleted || []);
            setBinaryHits(prev => prev.filter(hit => !deletedSet.has(hit.path)));
            setSelectedBinaryPaths(new Set());
            setDeleteSuccess(true);
            setDeleteModalOpen(false);
        } catch (err) {
            console.error('Delete failed:', err);
        } finally {
            setDeleting(false);
        }
    };

    const linkedWorktreesByRepoRaw = (() => {
        const byRepo = new Map<string, WorktreeHit[]>();
        for (const hit of worktreeHits) {
            const list = byRepo.get(hit.repoPath) || [];
            list.push(hit);
            byRepo.set(hit.repoPath, list);
        }
        return byRepo;
    })();

    const binariesByRepoRaw = (() => {
        const byRepo = new Map<string, BinaryHit[]>();
        for (const hit of binaryHits) {
            const list = byRepo.get(hit.repoPath) || [];
            list.push(hit);
            byRepo.set(hit.repoPath, list);
        }
        return byRepo;
    })();

    const { repos: filteredWorktreeRepos, linkedByRepo: filteredLinkedWorktreesByRepo } = filterWorktreeRepos(
        worktreeRepos,
        linkedWorktreesByRepoRaw,
        showWorktreeUnder10M,
    );
    const visibleWorktreeRepos = sortWorktreeRepos(filteredWorktreeRepos);

    const filteredBinariesByRepo = filterBinaryRepos(binariesByRepoRaw, showBinaryUnder1M);
    const visibleBinaryRepos = sortBinaryRepos(filteredBinariesByRepo);

    const namedByRepoRaw = (() => {
        const byRepo = new Map<string, NamedHit[]>();
        for (const hit of namedHits) {
            const list = byRepo.get(hit.repoPath) || [];
            list.push(hit);
            byRepo.set(hit.repoPath, list);
        }
        return byRepo;
    })();

    const vendorByRepoRaw = (() => {
        const byRepo = new Map<string, NamedHit[]>();
        for (const hit of vendorHits) {
            const list = byRepo.get(hit.repoPath) || [];
            list.push(hit);
            byRepo.set(hit.repoPath, list);
        }
        return byRepo;
    })();

    const filteredNamedByRepo = filterNamedRepos(namedByRepoRaw, showNamedUnder1M);
    const visibleNamedRepos = sortNamedRepos(filteredNamedByRepo);

    const filteredVendorByRepo = filterNamedRepos(vendorByRepoRaw, showVendorUnder1M);
    const visibleVendorRepos = sortNamedRepos(filteredVendorByRepo);

    const selectedNamedTotal = namedHits
        .filter(hit => selectedNamedPaths.has(hit.path))
        .reduce((sum, hit) => sum + hit.size, 0);

    const selectedVendorTotal = vendorHits
        .filter(hit => selectedVendorPaths.has(hit.path))
        .reduce((sum, hit) => sum + hit.size, 0);

    const coreCategories = new Set(['trash', 'temp', 'cache', 'log', 'swap']);
    const coreLocations = locations.filter(l => coreCategories.has(l.category));
    const softwareLocations = locations.filter(l => !coreCategories.has(l.category));
    const notDetectedSoftware = softwareLocations.filter(l => !l.detected);

    const renderCleanupPopover = (loc: TmpLocation) => {
        if (loc.category === 'swap') {
            return (
                <div data-testid={`cleanup-popover-${loc.category}`} style={{ maxWidth: '300px' }}>
                    <Text>Cannot be reclaimed -- managed by macOS</Text>
                </div>
            );
        }

        const suggestions = cleanupSuggestions[loc.category];
        if (!suggestions || suggestions.length === 0) {
            return (
                <div data-testid={`cleanup-popover-${loc.category}`} style={{ maxWidth: '300px' }}>
                    <Text type="secondary">No standard cleanup commands available. Consider manual removal of the listed paths.</Text>
                </div>
            );
        }

        return (
            <div data-testid={`cleanup-popover-${loc.category}`} style={{ maxWidth: '400px', maxHeight: '500px', overflowY: 'auto' }}>
                {suggestions.map((s, idx) => (
                    <div key={idx} data-testid={`cleanup-suggestion-${loc.category}-${idx}`} style={{ marginBottom: idx < suggestions.length - 1 ? '12px' : 0 }}>
                        <div>
                            <Text strong style={{ fontSize: '12px' }}>Command:</Text>
                            <div style={{ background: '#f5f5f5', padding: '4px 8px', borderRadius: '4px', marginTop: '2px' }}>
                                <Text code style={{ fontSize: '11px', wordBreak: 'break-all' }}>{s.command}</Text>
                            </div>
                        </div>
                        <div style={{ marginTop: '4px' }}>
                            <Text style={{ fontSize: '11px' }}>
                                <strong>Files removed:</strong> {s.removes}
                            </Text>
                        </div>
                        <div style={{ marginTop: '2px' }}>
                            <Text style={{ fontSize: '11px' }}>
                                <strong>Recoverability:</strong> {s.recoverable}
                            </Text>
                        </div>
                    </div>
                ))}
            </div>
        );
    };

    const renderCard = (loc: TmpLocation) => {
        const status = locStatuses[loc.label] || 'idle';
        const hasMultiBreakdown = loc.breakdownItems && loc.breakdownItems.length >= 2;
        const hasSinglePath = loc.breakdownItems && loc.breakdownItems.length === 1;

        return (
            <Col xs={24} sm={12} key={loc.label}>
                <Card
                    data-testid={`card-${loc.category}`}
                    size="small"
                    title={
                        <Space>
                            <Tag color={categoryColors[loc.category] || 'default'}>{loc.category}</Tag>
                            <span data-testid="card-label">{loc.label}</span>
                            <Popover
                                content={renderCleanupPopover(loc)}
                                trigger="click"
                            >
                                <span data-testid="cleanup-indicator" style={{ cursor: 'pointer', color: '#1677ff' }}>
                                    <QuestionCircleOutlined />
                                </span>
                            </Popover>
                            {status === 'pending' && (
                                <span data-testid="pending-badge"><ClockCircleOutlined style={{ color: '#faad14' }} /></span>
                            )}
                            {status === 'scanning' && (
                                <span data-testid="scanning-badge"><SyncOutlined spin style={{ color: '#1677ff' }} /></span>
                            )}
                            {status === 'done' && (
                                <span data-testid="done-badge"><CheckCircleOutlined style={{ color: '#52c41a' }} /></span>
                            )}
                        </Space>
                    }
                    extra={
                        <Space size={4}>
                            {loc.rebootSafe ? (
                                <Tag data-testid="reboot-safe-badge" color="green">Reboot Safe</Tag>
                            ) : (
                                <Tag data-testid="reboot-safe-badge" color="default">Cleared on Reboot</Tag>
                            )}
                            {loc.rebootSafe && !loc.reclaimable && (
                                <Tag data-testid="non-reclaimable-badge" color="orange">OS Managed</Tag>
                            )}
                        </Space>
                    }
                >
                    <div style={{ fontSize: '20px', fontWeight: 'bold' }} data-testid="card-size">
                        {formatBytes(loc.size)}
                    </div>
                    {loc.category === 'podman' && (
                        <Text type="secondary" style={{ fontSize: '11px', display: 'block' }}>On disk</Text>
                    )}
                    <Text type="secondary" style={{ fontSize: '12px' }}>
                        {loc.fileCount > 0 ? `${loc.fileCount} files` : '--'}
                    </Text>
                    {hasMultiBreakdown ? (
                        <div data-testid="breakdown-items" style={{ marginTop: '8px', paddingLeft: '8px', borderLeft: '2px solid #e8e8e8' }}>
                            {loc.breakdownItems!.map((item, idx) => (
                                <div key={idx} data-testid={`breakdown-row-${idx}`} style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '4px' }}>
                                    <Text data-testid={`breakdown-label-${idx}`} type="secondary" style={{ fontSize: '11px', fontFamily: 'monospace' }}>
                                        {item.path}
                                    </Text>
                                    <Text data-testid={`breakdown-size-${idx}`} strong style={{ fontSize: '13px' }}>
                                        {formatBytes(item.size)}
                                    </Text>
                                </div>
                            ))}
                        </div>
                    ) : (
                        hasSinglePath && loc.breakdownItems![0].path && (
                            <div style={{ marginTop: '4px' }}>
                                <Text data-testid="card-path" type="secondary" style={{ fontSize: '11px', fontFamily: 'monospace', wordBreak: 'break-all' }}>{loc.breakdownItems![0].path}</Text>
                            </div>
                        )
                    )}
                    {loc.vmInternal && loc.vmInternal.machineRunning && loc.vmInternal.items.length > 0 && (
                        <div data-testid="vm-internal-section" style={{ marginTop: '8px', paddingLeft: '8px', borderLeft: '2px solid #d3adf7' }}>
                            <Text type="secondary" style={{ fontSize: '11px', display: 'block', marginBottom: '4px' }}>Inside VM</Text>
                            {loc.vmInternal.items.map((item, idx) => (
                                <div key={idx} data-testid={`vm-internal-row-${idx}`} style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '4px' }}>
                                    <Text data-testid={`vm-internal-label-${idx}`} type="secondary" style={{ fontSize: '11px' }}>
                                        {item.label}
                                    </Text>
                                    <Text data-testid={`vm-internal-size-${idx}`} strong style={{ fontSize: '11px' }}>
                                        {formatBytes(item.size)}
                                    </Text>
                                </div>
                            ))}
                        </div>
                    )}
                    {loc.runtimeItems && loc.runtimeItems.length > 0 && (
                        <div data-testid="runtime-section" style={{ marginTop: '8px', paddingLeft: '8px', borderLeft: '2px solid #d9d9d9' }}>
                            <Text type="secondary" style={{ fontSize: '11px', display: 'block', marginBottom: '4px' }}>Runtime</Text>
                            {loc.runtimeItems.map((item, idx) => (
                                <div key={idx} data-testid={`runtime-row-${idx}`} style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '4px' }}>
                                    <Text data-testid={`runtime-label-${idx}`} type="secondary" style={{ fontSize: '11px' }}>
                                        {item.type}
                                    </Text>
                                    <Space size={4}>
                                        <Text data-testid={`runtime-count-${idx}`} style={{ fontSize: '11px' }}>
                                            {item.totalCount} {item.type === 'Images' ? 'images' : 'items'}
                                        </Text>
                                        <Text data-testid={`runtime-size-${idx}`} strong style={{ fontSize: '11px' }}>
                                            {formatBytes(item.size)}
                                        </Text>
                                    </Space>
                                </div>
                            ))}
                        </div>
                    )}
                </Card>
            </Col>
        );
    };

    return (
        <div style={{ padding: '24px', maxWidth: '900px', margin: '0 auto' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
                <h1 data-testid="page-heading" style={{ margin: 0 }}>Tmp Files Analyse</h1>
                <Space>
                    <Button
                        data-testid="start-scan-btn"
                        type="primary"
                        icon={<PlayCircleOutlined />}
                        onClick={startScan}
                        style={{ display: scanning ? 'none' : undefined }}
                    >
                        Start Scan
                    </Button>
                    <Button
                        data-testid="stop-scan-btn"
                        danger
                        icon={<StopOutlined />}
                        onClick={stopScan}
                        style={{ display: scanning ? undefined : 'none' }}
                    >
                        Stop Scan
                    </Button>
                </Space>
            </div>

            {platformError && (
                <Alert
                    message="Platform Not Supported"
                    description={platformError}
                    type="warning"
                    showIcon
                    style={{ marginBottom: '24px' }}
                />
            )}

            <Card
                data-testid="summary-bar"
                style={{ marginBottom: '24px', background: '#fafafa' }}
            >
                <Row gutter={24}>
                    <Col span={12}>
                        <Text type="secondary">Total Space</Text>
                        <div>
                            <Text strong style={{ fontSize: '18px' }} data-testid="total-size">
                                {formatBytes(totalSize)}
                            </Text>
                        </div>
                    </Col>
                    <Col span={12}>
                        <Text type="secondary">Reclaimable (safe to delete)</Text>
                        <div>
                            <Text strong style={{ fontSize: '18px', color: '#52c41a' }} data-testid="reclaimable-size">
                                {formatBytes(reclaimableSize)}
                            </Text>
                        </div>
                    </Col>
                </Row>
            </Card>

            <h2 data-testid="section-core-heading" style={{ fontSize: '16px', marginBottom: '12px' }}>System Locations</h2>
            <Row gutter={[16, 16]} style={{ marginBottom: '24px' }}>
                {coreLocations.map(renderCard)}
            </Row>

            <div data-testid="section-software">
                <h2 data-testid="section-software-heading" style={{ fontSize: '16px', marginBottom: '12px' }}>Developer Tools</h2>
                <Row gutter={[16, 16]} style={{ marginBottom: '24px' }}>
                    {softwareLocations.map(renderCard)}
                </Row>
            </div>

            {notDetectedSoftware.length > 0 && (
                <Collapse
                    data-testid="collapse-not-detected"
                    style={{ marginTop: '16px' }}
                    items={[{
                        key: 'not-detected',
                        label: 'Not Detected',
                        children: (
                            <div>
                                {notDetectedSoftware.map(loc => (
                                    <div key={loc.label} data-testid="not-detected-item" style={{ padding: '4px 0' }}>
                                        <Text data-testid="not-detected-item-name" type="secondary">{loc.label}</Text>
                                    </div>
                                ))}
                            </div>
                        ),
                    }]}
                />
            )}

            <h2 data-testid="section-repository-scans-heading" style={{ fontSize: '16px', marginTop: '24px', marginBottom: '12px' }}>
                Repository Scans
            </h2>

            <Card data-testid="worktrees-section" size="small" style={{ marginBottom: '16px' }}
                title={
                    <Space>
                        <span>Git Worktrees</span>
                        <span data-testid="worktrees-running-total" style={{ color: '#888' }}>
                            ({formatBytes(worktreeRunningTotal)})
                        </span>
                        {worktreesScanning && (
                            <span data-testid="worktrees-scanning-badge"><SyncOutlined spin style={{ color: '#1677ff' }} /></span>
                        )}
                        {worktreesDone && !worktreesScanning && (
                            <span data-testid="worktrees-done-badge"><CheckCircleOutlined style={{ color: '#52c41a' }} /></span>
                        )}
                    </Space>
                }
                extra={
                    <Space>
                        <Checkbox
                            data-testid="worktree-show-under-10m"
                            checked={showWorktreeUnder10M}
                            onChange={e => setShowWorktreeUnder10M(e.target.checked)}
                        >
                            &lt;10M
                        </Checkbox>
                        <Button
                            data-testid="worktrees-scan-btn"
                            size="small"
                            type="primary"
                            icon={<PlayCircleOutlined />}
                            onClick={startWorktreesScan}
                            style={{ display: worktreesScanning ? 'none' : undefined }}
                        >
                            Scan
                        </Button>
                        <Button
                            data-testid="worktrees-stop-btn"
                            size="small"
                            danger
                            icon={<StopOutlined />}
                            onClick={stopWorktreesScan}
                            disabled={!worktreesScanning}
                        >
                            Stop
                        </Button>
                    </Space>
                }
            >
                {visibleWorktreeRepos.length > 0 || worktreesScanning || worktreesDone ? (
                    <div data-testid="worktrees-tree" style={{ width: '100%', textAlign: 'left' }}>
                        {visibleWorktreeRepos.map(repo => {
                            const repoKey = testIdKey(repo.repoPath);
                            const linkedHits = sortLinkedWorktrees(filteredLinkedWorktreesByRepo.get(repo.repoPath) || []);
                            return (
                                <div key={repo.repoPath} style={{ marginBottom: '8px', textAlign: 'left' }}>
                                    <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '4px' }}>
                                        <span data-testid={`worktree-repo-row-${repoKey}`}>{repo.repoPath}</span>
                                        <Text strong data-testid={`worktree-repo-size-${repoKey}`}>{repo.sizeHuman || formatBytes(repo.size)}</Text>
                                    </div>
                                    {linkedHits.length > 0 && (
                                        <div style={{ paddingLeft: '16px', textAlign: 'left' }}>
                                            {linkedHits.map(hit => {
                                                const rowKey = testIdKey(hit.path);
                                                return (
                                                    <div key={hit.path} style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '4px', gap: 8 }}>
                                                        <span data-testid={`worktree-row-${rowKey}`}>{hit.path}</span>
                                                        <Text strong data-testid={`worktree-row-size-${rowKey}`}>{hit.sizeHuman || formatBytes(hit.size)}</Text>
                                                    </div>
                                                );
                                            })}
                                        </div>
                                    )}
                                </div>
                            );
                        })}
                    </div>
                ) : (
                    <Text type="secondary">Click Scan to discover git worktrees under ~</Text>
                )}
            </Card>

            <Card data-testid="binaries-section" size="small" style={{ marginBottom: '16px' }}
                title={
                    <Space>
                        <span>Binary files</span>
                        <span data-testid="binaries-running-total" style={{ color: '#888' }}>
                            ({formatBytes(binaryRunningTotal)})
                        </span>
                        {binariesScanning && (
                            <span data-testid="binaries-scanning-badge"><SyncOutlined spin style={{ color: '#1677ff' }} /></span>
                        )}
                        {binariesDone && !binariesScanning && (
                            <span data-testid="binaries-done-badge"><CheckCircleOutlined style={{ color: '#52c41a' }} /></span>
                        )}
                    </Space>
                }
                extra={
                    <Space>
                        <Checkbox
                            data-testid="binary-show-under-1m"
                            checked={showBinaryUnder1M}
                            onChange={e => setShowBinaryUnder1M(e.target.checked)}
                        >
                            &lt;1M
                        </Checkbox>
                        <Button
                            data-testid="binaries-scan-btn"
                            size="small"
                            type="primary"
                            icon={<PlayCircleOutlined />}
                            onClick={startBinariesScan}
                            style={{ display: binariesScanning ? 'none' : undefined }}
                        >
                            Scan
                        </Button>
                        <Button
                            data-testid="binaries-stop-btn"
                            size="small"
                            danger
                            icon={<StopOutlined />}
                            onClick={stopBinariesScan}
                            disabled={!binariesScanning}
                        >
                            Stop
                        </Button>
                    </Space>
                }
            >
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '12px' }}>
                    <Text data-testid="binary-selected-total">
                        Selected: {formatBytes(selectedBinaryTotal)} to clear
                    </Text>
                    <Button
                        data-testid="binary-delete-btn"
                        size="small"
                        danger
                        disabled={selectedBinaryPaths.size === 0}
                        onClick={() => setDeleteModalOpen(true)}
                    >
                        Delete Selected
                    </Button>
                </div>
                {deleteSuccess && (
                    <Alert data-testid="binary-delete-success" message="Selected binaries deleted" type="success" showIcon style={{ marginBottom: '12px' }} />
                )}
                {visibleBinaryRepos.length > 0 || binariesScanning || binariesDone ? (
                    <div data-testid="binaries-tree" style={{ width: '100%', textAlign: 'left' }}>
                        {visibleBinaryRepos.map(([repoPath, hits]) => {
                            const repoKey = testIdKey(repoPath);
                            const repoSize = hits.reduce((sum, h) => sum + h.size, 0);
                            const selectedInRepo = hits.filter(h => selectedBinaryPaths.has(h.path)).length;
                            const allSelected = selectedInRepo === hits.length;
                            const someSelected = selectedInRepo > 0 && !allSelected;
                            return (
                                <div key={repoPath} style={{ marginBottom: '8px', textAlign: 'left' }}>
                                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '4px' }}>
                                        <Checkbox
                                            data-testid={`binary-repo-checkbox-${repoKey}`}
                                            checked={allSelected}
                                            indeterminate={someSelected}
                                            onChange={e => toggleRepoBinaries(repoPath, e.target.checked)}
                                        >
                                            <span data-testid={`binary-repo-row-${repoKey}`}>{repoPath}</span>
                                        </Checkbox>
                                        <Text strong>{formatBytes(repoSize)}</Text>
                                    </div>
                                    <div style={{ paddingLeft: '24px', textAlign: 'left' }}>
                                        {hits.map(hit => {
                                            const rowKey = testIdKey(hit.path);
                                            return (
                                                <div key={hit.path} data-testid={`binary-row-${rowKey}`} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '4px' }}>
                                                    <Checkbox
                                                        data-testid={`binary-checkbox-${rowKey}`}
                                                        checked={selectedBinaryPaths.has(hit.path)}
                                                        onChange={e => toggleBinaryPath(hit.path, e.target.checked)}
                                                    >
                                                        <Space size={4}>
                                                            <span data-testid={`binary-path-${rowKey}`}>{hit.path}</span>
                                                            <Tag data-testid={`binary-kind-${rowKey}`}>{hit.kind}</Tag>
                                                        </Space>
                                                    </Checkbox>
                                                    <Text strong>{hit.sizeHuman}</Text>
                                                </div>
                                            );
                                        })}
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                ) : (
                    <Text type="secondary">Click Scan to find Go/Mach-O/ELF binaries in git repos</Text>
                )}
                {binariesDone && binaryHits.length === 0 && (
                    <Text data-testid="binaries-empty-state" type="secondary">No binaries found under ~</Text>
                )}
            </Card>

            <Card data-testid="node-modules-section" size="small" style={{ marginBottom: '16px' }}
                title={
                    <Space>
                        <span>node_modules</span>
                        <span data-testid="node-modules-running-total" style={{ color: '#888' }}>
                            ({formatBytes(namedRunningTotal)})
                        </span>
                        {namedScanning && (
                            <span data-testid="node-modules-scanning-badge"><SyncOutlined spin style={{ color: '#1677ff' }} /></span>
                        )}
                        {namedDone && !namedScanning && (
                            <span data-testid="node-modules-done-badge"><CheckCircleOutlined style={{ color: '#52c41a' }} /></span>
                        )}
                    </Space>
                }
                extra={
                    <Space>
                        <Checkbox
                            data-testid="node-modules-show-under-1m"
                            checked={showNamedUnder1M}
                            onChange={e => setShowNamedUnder1M(e.target.checked)}
                        >
                            &lt;1M
                        </Checkbox>
                        <Button
                            data-testid="node-modules-scan-btn"
                            size="small"
                            type="primary"
                            icon={<PlayCircleOutlined />}
                            onClick={startNamedScan}
                            style={{ display: namedScanning ? 'none' : undefined }}
                        >
                            Scan
                        </Button>
                        <Button
                            data-testid="node-modules-stop-btn"
                            size="small"
                            danger
                            icon={<StopOutlined />}
                            onClick={stopNamedScan}
                            disabled={!namedScanning}
                        >
                            Stop
                        </Button>
                    </Space>
                }
            >
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '12px' }}>
                    <Text data-testid="node-modules-selected-total">
                        Selected: {formatBytes(selectedNamedTotal)} to clear
                    </Text>
                    <Button
                        data-testid="node-modules-delete-btn"
                        size="small"
                        danger
                        disabled={selectedNamedPaths.size === 0}
                        onClick={() => setNamedDeleteModalOpen(true)}
                    >
                        Delete Selected
                    </Button>
                </div>
                {namedDeleteSuccess && (
                    <Alert data-testid="node-modules-delete-success" message="Selected node_modules deleted" type="success" showIcon style={{ marginBottom: '12px' }} />
                )}
                {visibleNamedRepos.length > 0 || namedScanning || namedDone ? (
                    <div data-testid="node-modules-tree" style={{ width: '100%', textAlign: 'left' }}>
                        {visibleNamedRepos.map(([repoPath, hits]) => {
                            const repoKey = testIdKey(repoPath);
                            const sortHits = sortNamedHits(hits);
                            const repoSize = sortHits.reduce((sum, h) => sum + h.size, 0);
                            const selectedInRepo = sortHits.filter(h => selectedNamedPaths.has(h.path)).length;
                            const allSelected = selectedInRepo === sortHits.length;
                            const someSelected = selectedInRepo > 0 && !allSelected;
                            return (
                                <div key={repoPath} style={{ marginBottom: '8px', textAlign: 'left' }}>
                                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '4px' }}>
                                        <Checkbox
                                            data-testid={`node-modules-repo-checkbox-${repoKey}`}
                                            checked={allSelected}
                                            indeterminate={someSelected}
                                            onChange={e => toggleNamedRepo(repoPath, e.target.checked)}
                                        >
                                            <span data-testid={`node-modules-repo-row-${repoKey}`}>{repoPath}</span>
                                        </Checkbox>
                                        <Text strong>{formatBytes(repoSize)}</Text>
                                    </div>
                                    <div style={{ paddingLeft: '24px', textAlign: 'left' }}>
                                        {sortHits.map(hit => {
                                            const rowKey = testIdKey(hit.path);
                                            return (
                                                <div key={hit.path} data-testid={`node-modules-row-${rowKey}`} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '4px' }}>
                                                    <Checkbox
                                                        data-testid={`node-modules-checkbox-${rowKey}`}
                                                        checked={selectedNamedPaths.has(hit.path)}
                                                        onChange={e => toggleNamedPath(hit.path, e.target.checked)}
                                                    >
                                                        <Space size={4}>
                                                            <span data-testid={`node-modules-path-${rowKey}`}>{hit.path}</span>
                                                            <span data-testid={`node-modules-name-${rowKey}`}>{hit.name}</span>
                                                            <span data-testid={`node-modules-repo-${rowKey}`}>{hit.repoName}</span>
                                                        </Space>
                                                    </Checkbox>
                                                    <Text strong data-testid={`node-modules-size-${rowKey}`}>{hit.sizeHuman}</Text>
                                                </div>
                                            );
                                        })}
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                ) : (
                    <Text type="secondary">Click Scan to find node_modules directories under ~</Text>
                )}
                {namedDone && namedHits.length === 0 && (
                    <Text data-testid="node-modules-empty-state" type="secondary">No node_modules found under ~</Text>
                )}
            </Card>

            <Modal
                data-testid="node-modules-delete-confirm-modal"
                title="Delete selected node_modules?"
                open={namedDeleteModalOpen}
                onCancel={() => setNamedDeleteModalOpen(false)}
                footer={[
                    <Button key="cancel" onClick={() => setNamedDeleteModalOpen(false)}>Cancel</Button>,
                    <Button
                        key="confirm"
                        data-testid="node-modules-delete-confirm-btn"
                        type="primary"
                        danger
                        loading={namedDeleting}
                        onClick={confirmDeleteNamed}
                    >
                        Delete
                    </Button>,
                ]}
            >
                <Text>
                    Delete {selectedNamedPaths.size} selected directory(s), freeing {formatBytes(selectedNamedTotal)}?
                </Text>
            </Modal>

            <Card data-testid="vendor-section" size="small" style={{ marginBottom: '16px' }}
                title={
                    <Space>
                        <span>vendor</span>
                        <span data-testid="vendor-running-total" style={{ color: '#888' }}>
                            ({formatBytes(vendorRunningTotal)})
                        </span>
                        {vendorScanning && (
                            <span data-testid="vendor-scanning-badge"><SyncOutlined spin style={{ color: '#1677ff' }} /></span>
                        )}
                        {vendorDone && !vendorScanning && (
                            <span data-testid="vendor-done-badge"><CheckCircleOutlined style={{ color: '#52c41a' }} /></span>
                        )}
                    </Space>
                }
                extra={
                    <Space>
                        <Checkbox
                            data-testid="vendor-show-under-1m"
                            checked={showVendorUnder1M}
                            onChange={e => setShowVendorUnder1M(e.target.checked)}
                        >
                            &lt;1M
                        </Checkbox>
                        <Button
                            data-testid="vendor-scan-btn"
                            size="small"
                            type="primary"
                            icon={<PlayCircleOutlined />}
                            onClick={startVendorScan}
                            style={{ display: vendorScanning ? 'none' : undefined }}
                        >
                            Scan
                        </Button>
                        <Button
                            data-testid="vendor-stop-btn"
                            size="small"
                            danger
                            icon={<StopOutlined />}
                            onClick={stopVendorScan}
                            disabled={!vendorScanning}
                        >
                            Stop
                        </Button>
                    </Space>
                }
            >
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '12px' }}>
                    <Text data-testid="vendor-selected-total">
                        Selected: {formatBytes(selectedVendorTotal)} to clear
                    </Text>
                    <Button
                        data-testid="vendor-delete-btn"
                        size="small"
                        danger
                        disabled={selectedVendorPaths.size === 0}
                        onClick={() => setVendorDeleteModalOpen(true)}
                    >
                        Delete Selected
                    </Button>
                </div>
                {vendorDeleteSuccess && (
                    <Alert data-testid="vendor-delete-success" message="Selected vendor directories deleted" type="success" showIcon style={{ marginBottom: '12px' }} />
                )}
                {visibleVendorRepos.length > 0 || vendorScanning || vendorDone ? (
                    <div data-testid="vendor-tree" style={{ width: '100%', textAlign: 'left' }}>
                        {visibleVendorRepos.map(([repoPath, hits]) => {
                            const repoKey = testIdKey(repoPath);
                            const sortHits = sortNamedHits(hits);
                            const repoSize = sortHits.reduce((sum, h) => sum + h.size, 0);
                            const selectedInRepo = sortHits.filter(h => selectedVendorPaths.has(h.path)).length;
                            const allSelected = selectedInRepo === sortHits.length;
                            const someSelected = selectedInRepo > 0 && !allSelected;
                            return (
                                <div key={repoPath} style={{ marginBottom: '8px', textAlign: 'left' }}>
                                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '4px' }}>
                                        <Checkbox
                                            data-testid={`vendor-repo-checkbox-${repoKey}`}
                                            checked={allSelected}
                                            indeterminate={someSelected}
                                            onChange={e => toggleVendorRepo(repoPath, e.target.checked)}
                                        >
                                            <span data-testid={`vendor-repo-row-${repoKey}`}>{repoPath}</span>
                                        </Checkbox>
                                        <Text strong>{formatBytes(repoSize)}</Text>
                                    </div>
                                    <div style={{ paddingLeft: '24px', textAlign: 'left' }}>
                                        {sortHits.map(hit => {
                                            const rowKey = testIdKey(hit.path);
                                            return (
                                                <div key={hit.path} data-testid={`vendor-row-${rowKey}`} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '4px' }}>
                                                    <Checkbox
                                                        data-testid={`vendor-checkbox-${rowKey}`}
                                                        checked={selectedVendorPaths.has(hit.path)}
                                                        onChange={e => toggleVendorPath(hit.path, e.target.checked)}
                                                    >
                                                        <Space size={4}>
                                                            <span data-testid={`vendor-path-${rowKey}`}>{hit.path}</span>
                                                            <span data-testid={`vendor-name-${rowKey}`}>{hit.name}</span>
                                                            <span data-testid={`vendor-repo-${rowKey}`}>{hit.repoName}</span>
                                                        </Space>
                                                    </Checkbox>
                                                    <Text strong data-testid={`vendor-size-${rowKey}`}>{hit.sizeHuman}</Text>
                                                </div>
                                            );
                                        })}
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                ) : (
                    <Text type="secondary">Click Scan to find vendor directories under ~</Text>
                )}
                {vendorDone && vendorHits.length === 0 && (
                    <Text data-testid="vendor-empty-state" type="secondary">No vendor directories found under ~</Text>
                )}
            </Card>

            <Modal
                data-testid="vendor-delete-confirm-modal"
                title="Delete selected vendor directories?"
                open={vendorDeleteModalOpen}
                onCancel={() => setVendorDeleteModalOpen(false)}
                footer={[
                    <Button key="cancel" onClick={() => setVendorDeleteModalOpen(false)}>Cancel</Button>,
                    <Button
                        key="confirm"
                        data-testid="vendor-delete-confirm-btn"
                        type="primary"
                        danger
                        loading={vendorDeleting}
                        onClick={confirmDeleteVendor}
                    >
                        Delete
                    </Button>,
                ]}
            >
                <Text>
                    Delete {selectedVendorPaths.size} selected directory(s), freeing {formatBytes(selectedVendorTotal)}?
                </Text>
            </Modal>

            <Modal
                data-testid="binary-delete-confirm-modal"
                title="Delete selected binaries?"
                open={deleteModalOpen}
                onCancel={() => setDeleteModalOpen(false)}
                footer={[
                    <Button key="cancel" onClick={() => setDeleteModalOpen(false)}>Cancel</Button>,
                    <Button
                        key="confirm"
                        data-testid="binary-delete-confirm-btn"
                        type="primary"
                        danger
                        loading={deleting}
                        onClick={confirmDeleteBinaries}
                    >
                        Delete
                    </Button>,
                ]}
            >
                <Text>
                    Delete {selectedBinaryPaths.size} selected file(s), freeing {formatBytes(selectedBinaryTotal)}?
                </Text>
            </Modal>
        </div>
    );
}

export default TmpFilesAnalyse;
