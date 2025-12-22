import { useState, useEffect, useRef } from 'react';
import { useSearchParams } from 'react-router-dom';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import { Pie } from 'react-chartjs-2';
import { Table, Input, Select, Button, Modal, Space, message, AutoComplete, Popover, Typography } from 'antd';
import { DeleteOutlined, ReloadOutlined, FolderOpenOutlined, FileOutlined, CopyOutlined, QuestionCircleOutlined, StopOutlined } from '@ant-design/icons';
import { DiskUsageAPI } from './api/disk_usage';
import type { FileInfo } from './api/disk_usage';

ChartJS.register(ArcElement, Tooltip, Legend);

function parseSize(sizeStr: string): number {
    const units = {
        'B': 1,
        'KB': 1024,
        'MB': 1024 * 1024,
        'GB': 1024 * 1024 * 1024,
        'TB': 1024 * 1024 * 1024 * 1024
    };
    const match = sizeStr.toUpperCase().match(/^([\d.]+)\s*([A-Z]+)$/);
    if (!match) return 0;
    const val = parseFloat(match[1]);
    const unit = match[2];
    return val * (units[unit as keyof typeof units] || 1);
}

function formatBytes(bytes: number, decimals = 2) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

interface FileNode extends FileInfo {
    children?: FileNode[];
    key: string; // Required for AntD Table
    path: string;
}

// Robust path join
function joinPath(p1: string, p2: string) {
    if (!p1) return p2;
    const sep = p1.includes('\\') ? '\\' : '/';
    if (p1 === '/') return '/' + p2;
    return p1.endsWith(sep) ? p1 + p2 : p1 + sep + p2;
}

// Helper to update a node deep in the tree
function updateTree(nodes: FileNode[], targetPath: string, updater: (node: FileNode) => FileNode): FileNode[] {
    return nodes.map(node => {
        if (node.path === targetPath) {
            return updater(node);
        }
        if (node.children && targetPath.startsWith(node.path + (node.path.endsWith('/') || node.path.endsWith('\\') ? '' : '/'))) {
            return { ...node, children: updateTree(node.children, targetPath, updater) };
        }
        return node;
    });
}

// Helper to sort items
function sortItems(items: FileNode[]): FileNode[] {
    return items.sort((a, b) => {
        if (a.isDir !== b.isDir) {
            return a.isDir ? -1 : 1;
        }
        return b.size - a.size;
    });
}

// Helper to remove a node from the tree
function removeNode(nodes: FileNode[], targetPath: string): FileNode[] {
    return nodes
        .filter(node => node.path !== targetPath)
        .map(node => {
            if (node.children) {
                return { ...node, children: removeNode(node.children, targetPath) };
            }
            return node;
        });
}

function getFileHelp(record: FileNode) {
    if (record.name.startsWith('podman-') && record.name.endsWith('.raw')) {
        return (
            <div style={{ maxWidth: 500, maxHeight: 600, overflowY: 'auto' }}>
                <Typography.Title level={5}>check info</Typography.Title>
                <Typography.Paragraph code copyable>
                    qemu-img info "{record.path}"
                </Typography.Paragraph>

                <Typography.Title level={5}>shrink</Typography.Title>
                <Typography.Paragraph code copyable>
                    qemu-img resize -f raw --shrink "{record.path}" 30G
                </Typography.Paragraph>

                <Typography.Title level={5}>podman management</Typography.Title>

                <div style={{ marginBottom: 8 }}>
                    <Typography.Text strong>inspect images</Typography.Text>
                    <Typography.Paragraph code copyable style={{ marginTop: 4 }}>
                        podman images -a
                    </Typography.Paragraph>
                </div>

                <div style={{ marginBottom: 8 }}>
                    <Typography.Text strong>remove image</Typography.Text>
                    <Typography.Paragraph code copyable style={{ marginTop: 4 }}>
                        podman rmi b510945a84a9
                    </Typography.Paragraph>
                </div>

                <div style={{ marginBottom: 8 }}>
                    <Typography.Text strong>check system size</Typography.Text>
                    <Typography.Paragraph code copyable style={{ marginTop: 4 }}>
                        podman system df -v
                    </Typography.Paragraph>
                </div>

                <div style={{ marginBottom: 8 }}>
                    <Typography.Text strong>reclaim storage</Typography.Text>
                    <Typography.Paragraph code copyable style={{ marginTop: 4 }}>
                        podman system prune -a -f --volumes
                    </Typography.Paragraph>
                </div>
            </div>
        );
    }
    return null;
}

export default function DiskUsage() {
    const [searchParams, setSearchParams] = useSearchParams();
    const currentUrlPath = searchParams.get('path');

    const [rootPath, setRootPath] = useState<string>(currentUrlPath || '');
    const [rootItems, setRootItems] = useState<FileNode[]>([]);
    const [loading, setLoading] = useState(false);
    const [view, setView] = useState<'list' | 'pie'>('list');
    const activeSources = useRef<Map<string, EventSource>>(new Map());
    const [filterSizeStr, setFilterSizeStr] = useState('1MB');
    const [minSizeBytes, setMinSizeBytes] = useState(1024 * 1024);

    useEffect(() => {
        const targetPath = currentUrlPath || '';
        setRootPath(targetPath);
        fetchUsage(targetPath, true);
    }, [currentUrlPath]);

    useEffect(() => {
        return () => {
            activeSources.current.forEach(es => es.close());
            activeSources.current.clear();
        };
    }, []);

    useEffect(() => {
        const bytes = parseSize(filterSizeStr);
        setMinSizeBytes(bytes);
    }, [filterSizeStr]);

    const fetchUsage = (dirPath: string, isRoot: boolean = false) => {
        if (activeSources.current.has(dirPath)) {
            activeSources.current.get(dirPath)?.close();
        }

        let effectivePath = dirPath;

        if (isRoot) {
            // Clear all active sources when switching root to avoid ghost updates
            activeSources.current.forEach(es => es.close());
            activeSources.current.clear();

            setLoading(true);
            setRootItems([]);
        }

        const es = DiskUsageAPI.streamUsage(dirPath, {
            onPath: (p) => {
                if (isRoot) setRootPath(p);
                effectivePath = p;
            },
            onItem: (item) => {
                const itemPath = joinPath(effectivePath, item.name);
                const nodeWithInfo: FileNode = {
                    ...item,
                    path: itemPath,
                    key: itemPath,
                    // Initialize children for dirs to make them expandable
                    children: item.isDir ? [] : undefined
                };

                if (isRoot) {
                    setRootItems(prev => {
                        const idx = prev.findIndex(i => i.name === item.name);
                        let newItems;
                        if (idx >= 0) {
                            newItems = [...prev];
                            const existingChildren = prev[idx].children;
                            const hasContent = existingChildren && existingChildren.length > 0;

                            newItems[idx] = {
                                ...prev[idx],
                                ...nodeWithInfo,
                                children: hasContent ? existingChildren : nodeWithInfo.children
                            };
                        } else {
                            newItems = [...prev, nodeWithInfo];
                        }
                        return sortItems(newItems);
                    });
                } else {
                    setRootItems(prev => updateTree(prev, effectivePath, (node) => {
                        const children = node.children || [];
                        const idx = children.findIndex(i => i.name === item.name);
                        let newChildren;
                        if (idx >= 0) {
                            newChildren = [...children];
                            const existingChildren = children[idx].children;
                            const hasContent = existingChildren && existingChildren.length > 0;

                            newChildren[idx] = {
                                ...children[idx],
                                ...nodeWithInfo,
                                children: hasContent ? existingChildren : nodeWithInfo.children
                            };
                        } else {
                            newChildren = [...children, nodeWithInfo];
                        }
                        return { ...node, children: sortItems(newChildren) };
                    }));
                }
            },
            onDone: () => {
                if (isRoot) setLoading(false);
                es.close();
                activeSources.current.delete(dirPath);
            },
            onError: (err) => {
                console.error(err);
                if (isRoot) setLoading(false);
                message.error('Error fetching data: ' + err);
                es.close();
                activeSources.current.delete(dirPath);
            }
        });

        activeSources.current.set(dirPath, es);
    };

    const handleMoveToTrash = (record: FileNode) => {
        Modal.confirm({
            title: 'Confirm Move to Trash',
            content: `Are you sure you want to move "${record.name}" to trash?`,
            okText: 'Move to Trash',
            okType: 'danger',
            cancelText: 'Cancel',
            onOk: async () => {
                try {
                    await DiskUsageAPI.moveToTrash(record.path);
                    message.success('Moved to trash');
                    // Remove from local state immediately
                    setRootItems(prev => removeNode(prev, record.path));
                } catch (e) {
                    message.error('Failed to move to trash: ' + e);
                }
            },
        });
    };

    const handleStop = () => {
        activeSources.current.forEach(es => es.close());
        activeSources.current.clear();
        setLoading(false);
        message.info('Analysis stopped');
    };

    const onSearch = () => {
        setSearchParams(prev => {
            const newParams = new URLSearchParams(prev);
            if (rootPath) {
                newParams.set('path', rootPath);
            } else {
                newParams.delete('path');
            }
            return newParams;
        });
    };

    const columns = [
        {
            title: 'Name',
            dataIndex: 'name',
            key: 'name',
            render: (text: string, record: FileNode) => (
                <Space>
                    {record.isDir ? <FolderOpenOutlined /> : <FileOutlined />}
                    {text}
                </Space>
            ),
        },
        {
            title: 'Size',
            dataIndex: 'size',
            key: 'size',
            render: (size: number, record: FileNode) => (
                <span>
                    {formatBytes(size)}
                    {record.status === 'pending' && <span style={{ color: '#999', marginLeft: 8, fontSize: '0.8em' }}>(...)</span>}
                </span>
            ),
        },
        {
            title: 'Action',
            key: 'action',
            render: (_: any, record: FileNode) => {
                const helpContent = getFileHelp(record);
                return (
                    <Space>
                        {helpContent && (
                            <Popover content={helpContent} title="Help" trigger="click" placement="left">
                                <Button
                                    type="text"
                                    icon={<QuestionCircleOutlined />}
                                    title="Show Help"
                                    onClick={(e) => e.stopPropagation()}
                                />
                            </Popover>
                        )}
                        <Button
                            type="text"
                            icon={<CopyOutlined />}
                            onClick={(e) => {
                                e.stopPropagation();
                                navigator.clipboard.writeText(record.path);
                                message.success('Path copied');
                            }}
                            title="Copy Path"
                        />
                        <Button
                            type="text"
                            danger
                            icon={<DeleteOutlined />}
                            onClick={(e) => { e.stopPropagation(); handleMoveToTrash(record); }}
                            title="Move to Trash"
                        />
                    </Space>
                );
            },
            width: 150,
            align: 'center' as const,
        },
    ];

    // Filter logic for recursive tree
    const filterNodes = (nodes: FileNode[]): FileNode[] => {
        return nodes
            .filter(node => node.size >= minSizeBytes || node.status === 'pending')
            .map(node => {
                if (node.children) {
                    return { ...node, children: filterNodes(node.children) };
                }
                return node;
            });
    };

    const filteredData = filterNodes(rootItems);

    return (
        <div style={{ padding: 24, height: '100vh', display: 'flex', flexDirection: 'column' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
                <Space.Compact style={{ flex: 1, marginRight: 16 }}>
                    <Input
                        value={rootPath}
                        onChange={e => setRootPath(e.target.value)}
                        onPressEnter={onSearch}
                        placeholder="Enter path..."
                        style={{ width: 'calc(100% - 100px)' }}
                    />
                    <Button
                        type="primary"
                        danger={loading}
                        onClick={loading ? handleStop : onSearch}
                        icon={loading ? <StopOutlined /> : <ReloadOutlined />}
                    >
                        {loading ? 'Stop' : 'Go'}
                    </Button>
                </Space.Compact>
                <Space>
                    <span>Filter &gt;</span>
                    <AutoComplete
                        value={filterSizeStr}
                        onChange={setFilterSizeStr}
                        options={[
                            { value: '1KB' },
                            { value: '1MB' },
                            { value: '10MB' },
                            { value: '100MB' },
                            { value: '1GB' },
                            { value: '10GB' },
                        ]}
                        style={{ width: 100 }}
                    />
                    <Select value={view} onChange={setView} options={[{ value: 'list', label: 'Tree View' }, { value: 'pie', label: 'Pie Chart' }]} />
                </Space>
            </div>

            <div style={{ flex: 1, overflow: 'auto' }}>
                {view === 'list' ? (
                    <Table
                        columns={columns}
                        dataSource={filteredData}
                        pagination={false}
                        expandable={{
                            onExpand: (expanded, record) => {
                                if (expanded && (!record.children || record.children.length === 0)) {
                                    fetchUsage(record.path);
                                }
                            }
                        }}
                        size="small"
                    />
                ) : (
                    <div style={{ display: 'flex', justifyContent: 'center', height: '100%' }}>
                        <PieChart items={filteredData} />
                    </div>
                )}
            </div>
        </div>
    );
}

function PieChart({ items }: { items: FileInfo[] }) {
    if (items.length === 0) return <div>No data (filtered)</div>;

    const total = items.reduce((acc, item) => acc + item.size, 0);
    if (total === 0) return <div>No data</div>;

    const topItems = items.slice(0, 10);
    const others = items.slice(10);
    const othersSize = others.reduce((acc, item) => acc + item.size, 0);

    const chartData = [...topItems];
    if (othersSize > 0) {
        chartData.push({ name: 'Others', size: othersSize, isDir: true, status: 'done' });
    }

    const data = {
        labels: chartData.map(item => item.name),
        datasets: [
            {
                label: 'Disk Usage',
                data: chartData.map(item => item.size),
                backgroundColor: [
                    '#FF6384',
                    '#36A2EB',
                    '#FFCE56',
                    '#4BC0C0',
                    '#9966FF',
                    '#FF9F40',
                    '#E7E9ED',
                    '#71B37C',
                ],
                borderColor: '#fff',
                borderWidth: 1,
            },
        ],
    };

    const options = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                position: 'right' as const,
            },
            tooltip: {
                callbacks: {
                    label: function (context: any) {
                        let label = context.label || '';
                        if (label) {
                            label += ': ';
                        }
                        if (context.parsed !== null) {
                            label += formatBytes(context.parsed);
                        }
                        return label;
                    }
                }
            }
        },
    };

    return <Pie data={data} options={options} />;
}
