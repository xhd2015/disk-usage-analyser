import { useState, useRef, useEffect } from 'react';
import { Button, Card, Space, Typography, Row, Col, Tag, Collapse, Alert } from 'antd';
import { PlayCircleOutlined, StopOutlined, SyncOutlined, CheckCircleOutlined, ClockCircleOutlined } from '@ant-design/icons';

const { Text } = Typography;

interface TmpLocation {
    path: string;
    label: string;
    category: string;
    size: number;
    fileCount: number;
    rebootSafe: boolean;
    detected: boolean;
    extraPaths?: string[];
    extraSizes?: number[];
    extraCounts?: number[];
}

function formatBytes(bytes: number): string {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}

const categoryColors: Record<string, string> = {
    trash: 'red',
    temp: 'orange',
    cache: 'blue',
    log: 'purple',
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
                p.label === loc.label ? { ...p, size: loc.size, fileCount: loc.fileCount, path: loc.path, extraSizes: loc.extraSizes, extraCounts: loc.extraCounts } : p
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

    const coreCategories = new Set(['trash', 'temp', 'cache', 'log']);
    const coreLocations = locations.filter(l => coreCategories.has(l.category));
    const softwareLocations = locations.filter(l => !coreCategories.has(l.category));
    const notDetectedSoftware = softwareLocations.filter(l => !l.detected);

    const renderCard = (loc: TmpLocation) => {
        const status = locStatuses[loc.label] || 'idle';

        return (
            <Col xs={24} sm={12} key={loc.label}>
                <Card
                    data-testid={`card-${loc.category}`}
                    size="small"
                    title={
                        <Space>
                            <Tag color={categoryColors[loc.category] || 'default'}>{loc.category}</Tag>
                            <span data-testid="card-label">{loc.label}</span>
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
                        loc.rebootSafe ? (
                            <Tag data-testid="reboot-safe-badge" color="green">Reboot Safe</Tag>
                        ) : (
                            <Tag data-testid="reboot-safe-badge" color="default">Cleared on Reboot</Tag>
                        )
                    }
                >
                    <div style={{ fontSize: '20px', fontWeight: 'bold' }} data-testid="card-size">
                        {formatBytes(loc.size)}
                    </div>
                    <Text type="secondary" style={{ fontSize: '12px' }}>
                        {loc.fileCount > 0 ? `${loc.fileCount} files` : '--'}
                    </Text>
                    {loc.path && (
                        <div style={{ marginTop: '4px' }}>
                            <Text data-testid="card-path" type="secondary" style={{ fontSize: '11px', fontFamily: 'monospace', wordBreak: 'break-all' }}>{loc.path}</Text>
                        </div>
                    )}
                    {loc.extraPaths && loc.extraPaths.length > 0 && (
                        <div data-testid="extra-breakdown" style={{ marginTop: '8px', paddingLeft: '8px', borderLeft: '2px solid #e8e8e8' }}>
                            {loc.extraPaths.map((ep, idx) => (
                                <div key={idx} data-testid={`extra-breakdown-row-${idx}`} style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '4px' }}>
                                    <Text data-testid={`extra-breakdown-label-${idx}`} type="secondary" style={{ fontSize: '11px', fontFamily: 'monospace' }}>
                                        {ep}
                                    </Text>
                                    <Text data-testid={`extra-breakdown-size-${idx}`} strong style={{ fontSize: '13px' }}>
                                        {formatBytes((loc.extraSizes && loc.extraSizes[idx]) ? loc.extraSizes[idx] : 0)}
                                    </Text>
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
        </div>
    );
}

export default TmpFilesAnalyse;
