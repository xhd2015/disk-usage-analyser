import { useState, useRef } from 'react';
import { Button, Card, Space, Typography, Row, Col, Tag } from 'antd';
import { PlayCircleOutlined, StopOutlined, SyncOutlined, CheckCircleOutlined, ClockCircleOutlined } from '@ant-design/icons';

const { Text } = Typography;

interface TmpLocation {
    path: string;
    label: string;
    category: string;
    size: number;
    fileCount: number;
    rebootSafe: boolean;
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
};

const defaultLocations: TmpLocation[] = [
    { path: '', label: 'User Trash', category: 'trash', size: 0, fileCount: 0, rebootSafe: true },
    { path: '', label: 'User Caches', category: 'cache', size: 0, fileCount: 0, rebootSafe: true },
    { path: '', label: 'User Logs', category: 'log', size: 0, fileCount: 0, rebootSafe: true },
    { path: '', label: 'System Temp', category: 'temp', size: 0, fileCount: 0, rebootSafe: false },
    { path: '', label: 'System Tmp', category: 'temp', size: 0, fileCount: 0, rebootSafe: false },
];

type LocStatus = 'idle' | 'pending' | 'scanning' | 'done';

function TmpFilesAnalyse() {
    const [scanning, setScanning] = useState(false);
    const [locations, setLocations] = useState<TmpLocation[]>(defaultLocations);
    const [locStatuses, setLocStatuses] = useState<Record<string, LocStatus>>({});
    const eventSourceRef = useRef<EventSource | null>(null);
    const [totalSize, setTotalSize] = useState(0);
    const [reclaimableSize, setReclaimableSize] = useState(0);

    const startScan = () => {
        setScanning(true);
        const pending: Record<string, LocStatus> = {};
        defaultLocations.forEach(l => { pending[l.label] = 'pending'; });
        setLocStatuses(pending);

        const es = new EventSource('/api/tmp-analyse');
        eventSourceRef.current = es;

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
                p.label === loc.label ? { ...p, size: loc.size, fileCount: loc.fileCount, path: loc.path } : p
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

            <Row gutter={[16, 16]}>
                {locations.map((loc) => {
                    const status = locStatuses[loc.label] || 'idle';
                    return (
                        <Col xs={24} sm={12} key={loc.label}>
                            <Card
                                data-testid={`card-${loc.category}`}
                                size="small"
                                title={
                                    <Space>
                                        <Tag color={categoryColors[loc.category]}>{loc.category}</Tag>
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
                            </Card>
                        </Col>
                    );
                })}
            </Row>
        </div>
    );
}

export default TmpFilesAnalyse;
