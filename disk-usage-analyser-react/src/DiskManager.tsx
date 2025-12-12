import { useEffect, useState } from 'react';
import { Table, Button, Space, message, Tag, Tooltip, Modal, Input } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { FolderOpenOutlined } from '@ant-design/icons';

interface DiskInfo {
    deviceID: string;
    name: string;
    size: number;
    available: number;
    mountPoint: string;
    content: string;
    isInternal: boolean;
    status: string;
    children?: DiskInfo[];
    key?: string;
}

export default function DiskManager() {
    const [disks, setDisks] = useState<DiskInfo[]>([]);
    const [loading, setLoading] = useState(false);
    const [passwordModalOpen, setPasswordModalOpen] = useState(false);
    const [password, setPassword] = useState('');
    const [pendingMountDevice, setPendingMountDevice] = useState<string | null>(null);

    const processDisks = (data: any[]): DiskInfo[] => {
        return data.map((d: any) => ({
            ...d,
            key: d.deviceID,
            children: d.children ? processDisks(d.children) : undefined
        }));
    };

    const fetchDisks = async () => {
        setLoading(true);
        try {
            const res = await fetch('/api/disks/list');
            if (!res.ok) {
                throw new Error(`Failed to fetch disks: ${res.statusText}`);
            }
            const data = await res.json();
            setDisks(processDisks(data));
        } catch (err: any) {
            message.error(err.message);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchDisks();
    }, []);

    const handleMount = async (deviceID: string, pwd?: string) => {
        try {
            const res = await fetch('/api/disks/mount', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ deviceID, password: pwd })
            });

            if (res.status === 401) {
                setPendingMountDevice(deviceID);
                setPasswordModalOpen(true);
                return;
            }

            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || res.statusText);
            }
            message.success('Mounted successfully');
            setPasswordModalOpen(false);
            setPassword('');
            setPendingMountDevice(null);
            fetchDisks();
        } catch (err: any) {
            message.error(`Failed to mount: ${err.message}`);
        }
    };

    const handlePasswordSubmit = () => {
        if (pendingMountDevice) {
            handleMount(pendingMountDevice, password);
        }
    };

    const handleUnmount = async (deviceID: string) => {
        try {
            const res = await fetch(`/api/disks/unmount?deviceID=${deviceID}`, { method: 'POST' });
            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || res.statusText);
            }
            message.success('Unmounted successfully');
            fetchDisks();
        } catch (err: any) {
            message.error(`Failed to unmount: ${err.message}`);
        }
    };

    const handleOpen = async (path: string) => {
        try {
            const res = await fetch(`/api/disks/open?path=${encodeURIComponent(path)}`, { method: 'POST' });
            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || res.statusText);
            }
        } catch (err: any) {
            message.error(`Failed to open path: ${err.message}`);
        }
    };

    const formatSize = (bytes: number) => {
        if (bytes === 0) return '-';
        const units = ['B', 'KB', 'MB', 'GB', 'TB'];
        let size = bytes;
        let unitIndex = 0;
        while (size >= 1024 && unitIndex < units.length - 1) {
            size /= 1024;
            unitIndex++;
        }
        return `${size.toFixed(2)} ${units[unitIndex]}`;
    };

    const columns: ColumnsType<DiskInfo> = [
        {
            title: 'Device',
            dataIndex: 'deviceID',
            key: 'deviceID',
            render: (text, record) => (
                <Space>
                    {text}
                    {record.isInternal && <Tag color="blue">Internal</Tag>}
                </Space>
            )
        },
        {
            title: 'Name',
            dataIndex: 'name',
            key: 'name',
            render: (text) => text || '-',
        },
        {
            title: 'Size',
            dataIndex: 'size',
            key: 'size',
            render: (size) => formatSize(size),
        },
        {
            title: 'Available',
            dataIndex: 'available',
            key: 'available',
            render: (available) => available ? formatSize(available) : '-',
        },
        {
            title: 'Type',
            dataIndex: 'content',
            key: 'content',
        },
        {
            title: 'Mount Point',
            dataIndex: 'mountPoint',
            key: 'mountPoint',
            render: (text) => (
                <Space>
                    {text || '-'}
                    {text && (
                        <Tooltip title="Open in Finder">
                            <Button
                                type="text"
                                icon={<FolderOpenOutlined />}
                                onClick={() => handleOpen(text)}
                                size="small"
                            />
                        </Tooltip>
                    )}
                </Space>
            ),
        },
        {
            title: 'Status',
            dataIndex: 'status',
            key: 'status',
            render: (status) => status ? <Tag color="orange">{status}</Tag> : '-',
        },
        {
            title: 'Actions',
            key: 'actions',
            render: (_, record) => (
                <Space size="middle">
                    {record.mountPoint ? (
                        <Button onClick={() => handleUnmount(record.deviceID)} danger>
                            Unmount
                        </Button>
                    ) : (
                        <Button
                            onClick={() => handleMount(record.deviceID)}
                            type="primary"
                            disabled={!!record.status} // Disable if status is set (e.g. Checking)
                        >
                            Mount
                        </Button>
                    )}
                </Space>
            ),
        },
    ];

    return (
        <div style={{ padding: '20px' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
                <h2>Disk Manager</h2>
                <Button onClick={fetchDisks} loading={loading}>
                    Refresh
                </Button>
            </div>

            <Table
                columns={columns}
                dataSource={disks}
                loading={loading}
                pagination={false}
                rowKey="deviceID"
                expandable={{
                    defaultExpandAllRows: true
                }}
            />

            <Modal
                title="Sudo Password Required"
                open={passwordModalOpen}
                onOk={handlePasswordSubmit}
                onCancel={() => {
                    setPasswordModalOpen(false);
                    setPendingMountDevice(null);
                    setPassword('');
                }}
            >
                <p>Please enter your sudo password to mount this disk:</p>
                <Input.Password
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    onPressEnter={handlePasswordSubmit}
                />
            </Modal>
        </div>
    );
}
