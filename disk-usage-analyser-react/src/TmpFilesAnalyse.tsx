import { useState, useRef, useEffect } from 'react';
import { Button, Card, Space, Typography, Row, Col, Tag, Collapse, Alert, Popover } from 'antd';
import { PlayCircleOutlined, StopOutlined, SyncOutlined, CheckCircleOutlined, ClockCircleOutlined, QuestionCircleOutlined } from '@ant-design/icons';

const { Text } = Typography;

interface TmpLocation {
    label: string;
    category: string;
    size: number;
    fileCount: number;
    rebootSafe: boolean;
    reclaimable: boolean;
    detected: boolean;
    breakdownItems?: { path: string; size: number; fileCount: number }[];
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
        { command: 'podman system prune', removes: 'Unused images, build cache', recoverable: 'Yes, re-pulled/rebuilt when needed' },
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

function TmpFilesAnalyse() {
    const [scanning, setScanning] = useState(false);
    const [locations, setLocations] = useState<TmpLocation[]>([]);
    const [locStatuses, setLocStatuses] = useState<Record<string, LocStatus>>({});
    const eventSourceRef = useRef<EventSource | null>(null);
    const [totalSize, setTotalSize] = useState(0);
    const [reclaimableSize, setReclaimableSize] = useState(0);
    const [platformError, setPlatformError] = useState<string | null>(null);

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

    const startScan = () => {
        setScanning(true);
        setPlatformError(null);

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
            es.close();
            eventSourceRef.current = null;
            setScanning(false);
        });

        es.addEventListener('progress', (e) => {
            const p = JSON.parse((e as MessageEvent).data);
            setLocations(prev => prev.map(loc =>
                loc.label === p.label
                    ? { ...loc, size: p.size, fileCount: p.fileCount }
                    : loc
            ));
            setLocStatuses(prev => ({ ...prev, [p.label]: 'scanning' }));
            setTotalSize(p.totalSize);
            setReclaimableSize(p.reclaimableSize);
        });

        es.addEventListener('location', (e) => {
            const loc: TmpLocation = JSON.parse((e as MessageEvent).data);
            setLocations(prev => prev.map(p =>
                p.label === loc.label ? { ...p, size: loc.size, fileCount: loc.fileCount, breakdownItems: loc.breakdownItems } : p
            ));
            setLocStatuses(prev => ({ ...prev, [loc.label]: 'done' }));
        });

        es.addEventListener('done', () => {
            es.close();
            eventSourceRef.current = null;
            setScanning(false);
        });

        es.addEventListener('server_error', (e) => {
            const d = JSON.parse((e as MessageEvent).data);
            console.error('SSE error:', d.error);
            es.close();
            eventSourceRef.current = null;
            setScanning(false);
        });

        es.onerror = () => {
            if (es.readyState === EventSource.CLOSED) return;
            es.close();
            eventSourceRef.current = null;
            setScanning(false);
        };
    };

    const stopScan = () => {
        if (eventSourceRef.current) {
            eventSourceRef.current.close();
            eventSourceRef.current = null;
        }
        setScanning(false);
    };

    const coreCategories = new Set(['trash', 'temp', 'cache', 'log', 'swap']);
    const coreLocations = locations.filter(l => coreCategories.has(l.category));
    const softwareLocations = locations.filter(l => !coreCategories.has(l.category));
    const notDetectedSoftware = softwareLocations.filter(l => !l.detected);

    const renderCleanupPopover = (loc: TmpLocation) => {
        if (loc.category === 'swap') {
            return (
                <div style={{ maxWidth: '300px' }}>
                    <Text>Cannot be reclaimed -- managed by macOS</Text>
                </div>
            );
        }

        const suggestions = cleanupSuggestions[loc.category];
        if (!suggestions || suggestions.length === 0) {
            return (
                <div style={{ maxWidth: '300px' }}>
                    <Text type="secondary">No standard cleanup commands available. Consider manual removal of the listed paths.</Text>
                </div>
            );
        }

        return (
            <div style={{ maxWidth: '400px' }}>
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
                                data-testid={`cleanup-popover-${loc.category}`}
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
        </div>
    );
}

export default TmpFilesAnalyse;
